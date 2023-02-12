package web

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"k8s-utils/auth"
	"k8s-utils/conf"
	"k8s-utils/connection"
	"k8s-utils/utils"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"log"

	//"github.com/joho/godotenv"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"net/http"

	"github.com/gin-gonic/gin"
)


func StartWebServer() {

	config := conf.LoadConfig()

	router := gin.Default()

	router.GET("/logs", func(c *gin.Context) {
		out := ""

		k8s := c.DefaultQuery("k8s", config.K8S.Cluster)
		namespace := c.DefaultQuery("namespace", config.K8S.Namespace)
        app := c.DefaultQuery("app", config.App.Name)
		timeSince := c.DefaultQuery("since", "240h")

		conn := connection.KubeConnection(k8s, false)
		appReq, _ := labels.NewRequirement("app.kubernetes.io/name", selection.Equals, []string{app}) 
	
		selector := labels.NewSelector()
		selector = selector.Add(*appReq)
		
		podList, err := conn.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		
		if err != nil {
			//log.Fatal("Error get list pods")
			log.Println("Error get list pods")
		}
		
		for _, pod := range podList.Items {
			for _ , container := range pod.Spec.Containers {
				containerLog := utils.GetPodLogs(conn, pod, container.Name, timeSince)

				out += strings.Repeat("=",100)+"\n"
				out += fmt.Sprintf("pod: %s container: %s since: %s\n", pod.Name,container.Name, timeSince)
				out += strings.Repeat("=",100)+"\n"
				out += containerLog+"\n"

			}
		}
		c.String(http.StatusOK, out)
	})

	router.GET("/health", func(c *gin.Context) {
			app := c.DefaultQuery("app", "logger")
			server := c.DefaultQuery("server", "smart")
			token := auth.GetToken(config)
			url := fmt.Sprintf("http://%s-app01-XXXXXXXXXXX.local:8085/%s/actuator/health", server, app)
			actuatorHealth := connection.GetActuatorHealth(url, token)
			actuatorHealthMarshal, _ := json.Marshal(actuatorHealth)
			c.String(http.StatusOK, string(actuatorHealthMarshal))

		})
	router.Run(":3000")
}
