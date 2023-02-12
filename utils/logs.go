package utils

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

func GetPodLogs(conn *kubernetes.Clientset, pod corev1.Pod, container string, since string) string {
	timeSince, _ := time.ParseDuration(since) 
	t := int64(timeSince.Seconds())
	podLogOpts := corev1.PodLogOptions{
		Container: container,
		SinceSeconds: &t,
	}
	req := conn.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		log.Fatal("error:", err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Fatal("error in copy information from podLogs to buf") 

	}

	defer podLogs.Close()

	return buf.String()
}