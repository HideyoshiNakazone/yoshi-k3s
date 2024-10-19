package cmd

import (
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/cluster"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/resources"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"gopkg.in/yaml.v3"
	"os"
)

type NodePair[T resources.NodeConfig] struct {
	Config  *T
	Options *[]string
}

func ParseConfig(configPath string) *CusterConfig {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config CusterConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil
	}

	return &config
}

func ConfigureFromConfig(config *CusterConfig, kubeconfigPath *string) error {
	c := createClusterFromConfig(config)

	var kubeconfigContent *[]byte

	for _, masterNode := range parseMasterNodes(config) {
		nodeConfig, err := c.ConfigureMasterNode(*masterNode.Config, *masterNode.Options)
		if err != nil {
			return err
		}

		if kubeconfigContent == nil {
			// When configuring a K3S cluster all certificates are the same,
			// so we can use the kubeconfig from any master node
			kubeconfigContent = nodeConfig
		}
	}

	for _, workerNode := range parseWorkerNodes(config) {
		err := c.ConfigureWorkerNode(*workerNode.Config, *workerNode.Options)
		if err != nil {
			return err
		}
	}

	if kubeconfigContent == nil {
		return fmt.Errorf("invalid KUBECONFIG Returned, check the cluster state")
	}

	err := writeKubeconfig(kubeconfigPath, *kubeconfigContent)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFromConfig(config *CusterConfig) error {
	c := createClusterFromConfig(config)

	for _, masterNode := range parseMasterNodes(config) {
		err := c.DestroyMasterNode(*masterNode.Config)
		if err != nil {
			return err
		}
	}

	for _, workerNode := range parseWorkerNodes(config) {
		err := c.DestroyWorkerNode(*workerNode.Config)
		if err != nil {
			return err
		}
	}

	return nil
}

func createClusterFromConfig(config *CusterConfig) *cluster.K3sCluster {
	if config.Cluster.Version == "" {
		return cluster.NewK3sClient(
			config.Cluster.Token,
			config.Cluster.ServerAddress,
		)
	}

	return cluster.NewK3sClientWithVersion(
		config.Cluster.Version,
		config.Cluster.Token,
		config.Cluster.ServerAddress,
	)
}

func parseMasterNodes(config *CusterConfig) []NodePair[resources.NodeConfig] {
	var masterNodes []NodePair[resources.NodeConfig]

	for _, node := range config.MasterNodes {
		masterConfig := resources.NewNodeConfig(
			node.Name,
			ssh_handler.NewSshConfig(
				node.Connection.Host,
				node.Connection.Port,
				node.Connection.User,
				node.Connection.Password,
				node.Connection.PrivateKey,
				node.Connection.PrivateKeyPassphrase,
			),
		)

		masterNodes = append(masterNodes, NodePair[resources.NodeConfig]{
			Config:  masterConfig,
			Options: &node.Options,
		})
	}

	return masterNodes
}

func parseWorkerNodes(config *CusterConfig) []NodePair[resources.NodeConfig] {
	var workerNodes []NodePair[resources.NodeConfig]

	for _, node := range config.WorkerNodes {
		workerConfig := resources.NewNodeConfig(
			node.Name,
			ssh_handler.NewSshConfig(
				node.Connection.Host,
				node.Connection.Port,
				node.Connection.User,
				node.Connection.Password,
				node.Connection.PrivateKey,
				node.Connection.PrivateKeyPassphrase,
			),
		)
		workerNodes = append(workerNodes, NodePair[resources.NodeConfig]{
			Config:  workerConfig,
			Options: &node.Options,
		})
	}

	return workerNodes
}

func writeKubeconfig(kubeconfigPath *string, kubeconfigContent []byte) error {
	return os.WriteFile(*kubeconfigPath, kubeconfigContent, 0644)
}
