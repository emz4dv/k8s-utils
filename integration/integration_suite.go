package integration

import (
	"k8s-utils/auth"
	"k8s-utils/conf"
	"k8s-utils/connection"
	"k8s-utils/utils"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/stretchr/testify/suite"

	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"


)
type intSpringbootSuite struct {
	suite.Suite
	config *conf.Config
	connection *kubernetes.Clientset
}

func (s *intSpringbootSuite) getPodsList() *v1.PodList{
	appReq, _ := labels.NewRequirement("app.kubernetes.io/name", selection.Equals, []string{s.config.App.Artifact}) 
	
	selector := labels.NewSelector()
	selector = selector.Add(*appReq)

	podList, err := s.connection.CoreV1().Pods(s.config.K8S.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		s.FailNowf("error getting list of pods", err.Error())
	} else if len(podList.Items) == 0 {
		s.FailNow("no pods found")
	}
	return podList
}

func (s *intSpringbootSuite) springbootStart() {
	config := s.config
	var containerLog string

	message, err := retry.DoWithRetryE(s.T(), 
		"Try to query logs", 
		config.Tests.Springboot.Timeouts.Counts, 
		time.Duration(config.Tests.Springboot.Timeouts.Seconds)*time.Second,
		func() (string, error) {
			matcher := `Started \w+ServerApp in (\d+\.\d+) seconds \(JVM running for (\d+\.\d+\))`
			re, err := regexp.Compile(matcher)
			if err != nil {
				return "Error compile regex" , err
			}
			for _, pod := range s.getPodsList().Items {
				for _ , container := range pod.Spec.Containers {
		
					containerLog = utils.GetPodLogs(s.connection, pod, container.Name, config.Tests.Springboot.Since)
					res := re.FindStringSubmatch(containerLog)	

					if res != nil {
						s.T().Log(res[0])
					} else {
						return "" , errors.New("application launch string not found in log")
					}
				}
			}
			return "Application started successfully!", nil
		})
		
	if err != nil {
		s.T().Log(containerLog)
		s.FailNowf("Application run failed!", err.Error())
	} else {
		s.T().Log(message)
	}
}


func (s *intSpringbootSuite) actuatorHealth() {
	token := auth.GetToken(s.config)

	// временное решение
	url := fmt.Sprintf("http://%s-XXXXXXXXXXXXXXXXX.local:8085/%s/actuator/health",s.config.Auth.Realm, s.config.App.Name)

	actuatorHealth := connection.GetActuatorHealth(url, token)
	statusMsg := fmt.Sprintf("Actuator status: %s", actuatorHealth.Status)
	if actuatorHealth.Status == "UP" {
		s.T().Log(statusMsg)
	} else {
		s.FailNow(statusMsg)
	}
}


func (s *intSpringbootSuite) swagger() {
	token := auth.GetToken(s.config)

	// временное решение
	url := fmt.Sprintf("http://%s-XXXXXXXXXXXXXXX.local:8085/%s/v2/api-docs",s.config.Auth.Realm, s.config.App.Name)
	// запуск в k8s
	//url := fmt.Sprintf("http://%s:%s/%s/v2/api-docs",s.config.Test.Auth.Realm, s.config.Test.App.Name)
	appSwagger := connection.GetSwagger(url, token)
	if appSwagger.Paths == nil {
		s.FailNow("No paths in swagger")

	} else {
		s.T().Logf("Swagger version: %s", appSwagger.Swagger)
	}

}

func (s *intSpringbootSuite) podStatus() {

	for _, pod := range s.getPodsList().Items {
		message, err := retry.DoWithRetryE(s.T(), 
		"Try to check pod status", 
		s.config.Tests.PodStatus.Timeouts.Counts, 
		time.Duration(s.config.Tests.PodStatus.Timeouts.Seconds)*time.Second,
		func() (string, error) {

			getContainerStatusMessage := func(p v1.Pod) {
				for _, cont := range pod.Status.ContainerStatuses {
					if !cont.Ready {
						s.T().Error(cont.State.Waiting.Message)
					}
				}
			}

			switch pod.Status.Phase {
			case v1.PodRunning:
				return fmt.Sprintf("Pod status: %s", v1.PodRunning), nil
			case v1.PodFailed:
				getContainerStatusMessage(pod)
				return fmt.Sprintf("Pod status: %s", v1.PodFailed), errors.New("pod status Failed")
			case v1.PodPending:
				getContainerStatusMessage(pod)
				return fmt.Sprintf("Pod status: %s", v1.PodPending), errors.New("pod status Pending")
			}
			
			return "", errors.New("unknown pod status")
		})
		if err != nil {
			s.FailNow(message)
		} else {
			s.T().Logf("Start time: %s", pod.Status.StartTime)
			s.T().Log(message)
			
		}

	}
	

}