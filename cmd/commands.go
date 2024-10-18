package cmd

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/cluster"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/resources"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"gopkg.in/yaml.v3"
	"os"
)

type NodePair[T resources.NodeConfigInterface] struct {
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

	for _, masterNode := range parseMasterNodes(config) {
		err := c.ConfigureMasterNode(*masterNode.Config, *masterNode.Options)
		if err != nil {
			return err
		}
	}

	for _, workerNode := range parseWorkerNodes(config) {
		err := c.ConfigureWorkerNode(*workerNode.Config, *workerNode.Options)
		if err != nil {
			return err
		}
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
		return cluster.NewK3sClient(config.Cluster.Token)
	}

	return cluster.NewK3sClientWithVersion(config.Cluster.Version, config.Cluster.Token)
}

func parseMasterNodes(config *CusterConfig) []NodePair[resources.K3sMasterNodeConfig] {
	var masterNodes []NodePair[resources.K3sMasterNodeConfig]

	for _, node := range config.MasterNodes {
		masterConfig := resources.NewK3sMasterNodeConfig(
			ssh_handler.NewSshConfig(
				node.Connection.Host,
				node.Connection.Port,
				node.Connection.User,
				node.Connection.Password,
				node.Connection.PrivateKey,
				node.Connection.PrivateKeyPassphrase,
			),
		)

		masterNodes = append(masterNodes, NodePair[resources.K3sMasterNodeConfig]{
			Config:  masterConfig,
			Options: &node.Options,
		})
	}

	return masterNodes
}

func parseWorkerNodes(config *CusterConfig) []NodePair[resources.K3sWorkerNodeConfig] {
	var workerNodes []NodePair[resources.K3sWorkerNodeConfig]

	for _, node := range config.WorkerNodes {
		workerConfig := resources.NewK3sWorkerNodeConfig(
			node.ServerAddress,
			ssh_handler.NewSshConfig(
				node.Connection.Host,
				node.Connection.Port,
				node.Connection.User,
				node.Connection.Password,
				node.Connection.PrivateKey,
				node.Connection.PrivateKeyPassphrase,
			),
		)
		workerNodes = append(workerNodes, NodePair[resources.K3sWorkerNodeConfig]{
			Config:  workerConfig,
			Options: &node.Options,
		})
	}

	return workerNodes
}
