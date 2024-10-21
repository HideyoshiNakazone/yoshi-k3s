package kubeconfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func UpdateServerAddress(kubeconfigData *[]byte, serverAddress string) (*[]byte, error) {
	var kubeconfig = NewKubeconfigModel(kubeconfigData)
	if kubeconfig == nil {
		return nil, fmt.Errorf("invalid KUBECONFIG file, probably failed cluster configuration")
	}

	for i := range kubeconfig.Clusters {
		kubeconfig.Clusters[i].Cluster.Server = fmt.Sprintf("https://%s:6443", serverAddress)
	}

	newKubeconfigData, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return nil, err
	}

	return &newKubeconfigData, nil
}
