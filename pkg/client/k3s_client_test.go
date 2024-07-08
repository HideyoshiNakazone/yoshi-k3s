package client

import (
	"terraform-yoshi-k3s/pkg/resources"
	"terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sClient_ConfigureNode(t *testing.T) {
	k3sToken := "token"
	k3sVersion := "v1.30.2+k3s1"

	c := NewK3sClient()

	masterNodeArgs := []string{
		"--disable traefik",
		"--node-label node_type=master",
		"--snapshotter native",
	}

	var masterNodeConfig = resources.K3sMasterNodeConfig{
		Host:    "localhost",
		Token:   k3sToken,
		Version: k3sVersion,
		ConnectionConfig: ssh_handler.SshConfig{
			Host:     "localhost",
			Port:     "2222",
			User:     "sshuser",
			Password: "password",
		},
	}

	err := c.ConfigureMasterNode(masterNodeConfig, masterNodeArgs)
	if err != nil {
		t.Errorf("Error configuring master node: %v", err)
		return
	}

	var workerNodeConfig = resources.K3sWorkerNodeConfig{
		Server:  "master_node",
		Host:    "localhost",
		Token:   k3sToken,
		Version: k3sVersion,
		ConnectionConfig: ssh_handler.SshConfig{
			Host:     "localhost",
			Port:     "3333",
			User:     "sshuser",
			Password: "password",
		},
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
