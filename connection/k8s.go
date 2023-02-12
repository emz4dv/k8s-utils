package connection

import (
	//"flag"

	"log"
	"path/filepath"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)
func buildConfigFromCtx(ctx string, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
			&clientcmd.ConfigOverrides{
					CurrentContext: ctx,
			}).ClientConfig()
}
func KubeConnection(ctx string, inCluster bool) *kubernetes.Clientset {
	var kubeconfig string
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			kubeconfig = ""
		}
	
		config, err = buildConfigFromCtx(ctx, kubeconfig)
	}

	if err != nil {
		log.Fatal("error:", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("error:", err)
	}

	return clientset
}