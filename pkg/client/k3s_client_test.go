package client

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/resources"
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sClient_ConfigureNode(t *testing.T) {
	k3sToken := "token"
	k3sVersion := "v1.30.2+k3s2"

	masterNodes := NodeMapping[resources.K3sMasterNodeConfig]{}
	workerNodes := NodeMapping[resources.K3sWorkerNodeConfig]{}

	c := NewK3sClient(
		masterNodes,
		workerNodes,
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
		"localhost",
		k3sToken,
		k3sVersion,
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
		"localhost",
		k3sToken,
		k3sVersion,
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

	if !c.IsMasterNodeConfigured(masterNodeConfig.GetHost()) {
		t.Errorf("Master node not added to client")
		return
	}
	t.Logf("Master nodes in the cluster: %v", len(c.masterNodes))

	if !c.IsWorkerNodeConfigured(workerNodeConfig.GetHost()) {
		t.Errorf("Worker node not added to client")
		return
	}
	t.Logf("Worker nodes in the cluster: %v", len(c.workerNodes))

	err = c.DestroyMasterNode(masterNodeConfig.GetHost())
	if err != nil {
		t.Errorf("Error destroying master node: %v", err)
		return
	}

	err = c.DestroyWorkerNode(workerNodeConfig.GetHost())
	if err != nil {
		t.Errorf("Error destroying worker node: %v", err)
		return
	}
}
