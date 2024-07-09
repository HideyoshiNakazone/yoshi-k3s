package resources

import (
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
)

type K3sWorkerNodeConfig struct {
	server string

	host    string
	token   string
	version string

	connectionConfig *ssh_handler.SshConfig
}

func NewK3sWorkerNodeConfig(server string, host string, token string, version string, connectionConfig *ssh_handler.SshConfig) *K3sWorkerNodeConfig {
	return &K3sWorkerNodeConfig{
		server:           server,
		host:             host,
		token:            token,
		version:          version,
		connectionConfig: connectionConfig,
	}
}

func (k K3sWorkerNodeConfig) GetServer() string {
	return k.server
}
func (k K3sWorkerNodeConfig) GetHost() string {
	return k.host
}
func (k K3sWorkerNodeConfig) GetToken() string {
	return k.token
}
func (k K3sWorkerNodeConfig) GetVersion() string {
	return k.version
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
func (k K3sWorkerNodeConfig) HasChanged(other NodeConfigInterface) bool {
	return hasChanged(k, other) || k.GetServer() != other.(*K3sWorkerNodeConfig).GetServer()
}
