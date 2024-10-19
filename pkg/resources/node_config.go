package resources

import (
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
)

type NodeConfig struct {
	name             string
	connectionConfig *ssh_handler.SshConfig
}

func NewNodeConfig(name string, connectionConfig *ssh_handler.SshConfig) *NodeConfig {
	return &NodeConfig{
		name:             name,
		connectionConfig: connectionConfig,
	}
}

func (n NodeConfig) GetName() *string {
	return &n.name
}

func (n NodeConfig) GetConnectionConfig() *ssh_handler.SshConfig {
	return n.connectionConfig
}

func (n NodeConfig) IsValid() error {
	var err error

	err = n.GetConnectionConfig().IsValid()
	if err != nil {
		return err
	}

	if n.GetName() == nil || *n.GetName() == "" {
		return fmt.Errorf("name is empty")
	}

	return nil
}
