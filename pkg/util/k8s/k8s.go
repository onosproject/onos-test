package k8s

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetClient returns a Kubernetes client
func GetClient() (client.Client, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	return client.New(config, client.Options{})
}

// GetClientset returns a Kubernetes clientset
func GetClientset() (*kubernetes.Clientset, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// GetRestConfig returns the Kubernetes REST API configuration
func GetRestConfig() (*rest.Config, error) {
	restconfig, err := rest.InClusterConfig()
	if err == nil {
		return restconfig, nil
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home := getHomeDir()
		if home == "" {
			return nil, errors.New("no home directory configured")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// getHomeDir returns the user's home directory if defined by environment variables
func getHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
