package kubeconfig

import (
	"fmt"
	"testing"
)

var exampleKubeconfigData = []byte(`
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: empty-cert
    server: 127.0.0.1:6443
  name: default
contexts:
- context:
    cluster: default
    user: default
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: default
  user:
    password: empty-password
    username: admin
`)

func TestUpdateServerAddress(t *testing.T) {
	newServerAddress := "new-server-address"
	newClusterConfig, err := UpdateServerAddress(&exampleKubeconfigData, newServerAddress)
	if err != nil || newClusterConfig == nil {
		t.Errorf("Error updating server address: %v", err)
		return
	}

	newClusterConfigModel := NewKubeconfigModel(newClusterConfig)
	if newClusterConfigModel == nil {
		t.Errorf("Error creating new kubeconfig model")
		return
	}

	addressString := fmt.Sprintf("%s:6443", newServerAddress)
	if newClusterConfigModel.Clusters[0].Cluster.Server != addressString {
		t.Errorf("Error updating server address")
		return
	}
}
