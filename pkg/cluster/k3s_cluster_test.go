package cluster

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/resources"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sCluster_ConfigureNode(t *testing.T) {
	k3sToken := "token"
	k3sVersion := "v1.30.2+k3s2"

	c := NewK3sClientWithVersion(
		k3sVersion,
		k3sToken,
	)

	masterNodeArgs := []string{
		"--disable traefik",
		"--node-label node_type=master",
		"--snapshotter native",
	}

	var masterNodeSshConfig = ssh_handler.NewSshConfig(
		"localhost",
		"2222",
		"sshuser",
		"password",
		"",
		"",
	)
	var masterNodeConfig = resources.NewK3sMasterNodeConfig(
		masterNodeSshConfig,
	)

	err := c.ConfigureMasterNode(*masterNodeConfig, masterNodeArgs)
	if err != nil {
		t.Errorf("Error configuring master node: %v", err)
		return
	}

	var workerNodeSshConfig = ssh_handler.NewSshConfig(
		"localhost",
		"3333",
		"sshuser",
		"password",
		"",
		"",
	)
	var workerNodeConfig = resources.NewK3sWorkerNodeConfig(
		"master_node",
		workerNodeSshConfig,
	)

	workerNodeArgs := []string{
		"--node-label node_type=worker",
		"--snapshotter native",
	}

	err = c.ConfigureWorkerNode(*workerNodeConfig, workerNodeArgs)
	if err != nil {
		t.Errorf("Error configuring worker node: %v", err)
		return
	}

	err = c.DestroyMasterNode(*masterNodeConfig)
	if err != nil {
		t.Errorf("Error destroying master node: %v", err)
		return
	}

	err = c.DestroyWorkerNode(*workerNodeConfig)
	if err != nil {
		t.Errorf("Error destroying worker node: %v", err)
		return
	}
}
