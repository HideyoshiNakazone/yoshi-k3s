package client

import (
	"terraform-yoshi-k3s/pkg/resources"
	"terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sClient_ConfigureNode(t *testing.T) {
	host := "localhost"
	port := "3333"
	user := "sshuser"
	password := "password"

	c := NewK3sClient()

	masterNodeArgs := []string{
		"--disable traefik",
		"--node-label node_type=master",
		"--snapshotter native",
	}

	var masterNodeConfig = resources.K3sMasterNodeConfig{
		Host:    host,
		Token:   "token",
		Version: "v1.30.2+k3s1",
		ConnectionConfig: ssh_handler.SshConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
		},
	}

	err := c.ConfigureMasterNode(masterNodeConfig, masterNodeArgs)
	if err != nil {
		t.Errorf("Error configuring master node: %v", err)
		return
	}

	var workerNodeConfig = resources.K3sWorkerNodeConfig{
		Server:              "master_node",
		K3sMasterNodeConfig: masterNodeConfig,
	}

	workerNodeArgs := []string{
		"--node-label node_type=worker",
		"--snapshotter native",
	}

	err = c.ConfigureWorkerNode(workerNodeConfig, workerNodeArgs)
	if err != nil {
		t.Errorf("Error configuring worker node: %v", err)
		return
	}

	if len(c.masterNodes) == 0 {
		t.Errorf("Master node not added to client")
		return
	}
	t.Logf("Master nodes in the cluster: %v", len(c.masterNodes))

	if len(c.workerNodes) == 0 {
		t.Errorf("Worker node not added to client")
		return
	}
	t.Logf("Worker nodes in the cluster: %v", len(c.workerNodes))
}
