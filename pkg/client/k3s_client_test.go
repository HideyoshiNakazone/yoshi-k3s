package client

import (
	"terraform-yoshi-k3s/pkg/resources"
	"terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sClient_ConfigureMasterNode(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	nodeArgs := []string{
		"--disable traefik",
		"--node-label node_type=master",
		"--snapshotter native",
	}

	var nodeConfig = resources.K3sMasterNodeConfig{
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

	c := NewK3sClient()

	err := c.ConfigureMasterNode(nodeConfig, nodeArgs)
	if err != nil {
		t.Errorf("Error configuring master node: %v", err)
		return
	}
}

func TestK3sClient_ConfigureWorkerNode(t *testing.T) {
	host := "localhost"
	port := "3333"
	user := "sshuser"
	password := "password"

	nodeArgs := []string{
		"--node-label node_type=worker",
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

	var nodeConfig = resources.K3sWorkerNodeConfig{
		Server:              "master_node",
		K3sMasterNodeConfig: masterNodeConfig,
	}

	c := NewK3sClient()

	err := c.ConfigureWorkerNode(nodeConfig, nodeArgs)
	if err != nil {
		t.Errorf("Error configuring worker node: %v", err)
		return
	}
}
