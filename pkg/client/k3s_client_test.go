package client

import (
	"testing"
	"yoshi_k3s/pkg/resources"
	"yoshi_k3s/pkg/ssh_handler"
)

func TestK3sClient_ConfigureMasterNode(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	nodeArgs := []string{
		"--disable traefik",
		"--node-label node_type=master",
	}

	var nodeConfig = resources.K3sMasterNodeConfig{
		Host:    host,
		Token:   "token",
		Version: "latest",
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
		return
	}
}
