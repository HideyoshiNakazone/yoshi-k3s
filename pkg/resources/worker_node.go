package resources

import (
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
)

type K3sWorkerNodeConfig struct {
	server string

	connectionConfig *ssh_handler.SshConfig
}

func NewK3sWorkerNodeConfig(server string, connectionConfig *ssh_handler.SshConfig) *K3sWorkerNodeConfig {
	return &K3sWorkerNodeConfig{
		server:           server,
		connectionConfig: connectionConfig,
	}
}

func (k K3sWorkerNodeConfig) GetServer() string {
	return k.server
}
func (k K3sWorkerNodeConfig) GetConnectionConfig() *ssh_handler.SshConfig {
	return k.connectionConfig
}
func (k K3sWorkerNodeConfig) IsValid() error {
	if k.GetServer() == "" {
		return fmt.Errorf("server is empty")
	}
	return isNodeConfigValid(k)
}
