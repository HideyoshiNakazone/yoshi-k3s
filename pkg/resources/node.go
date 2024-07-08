package resources

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"fmt"
)

type NodeConfigInterface interface {
	GetHost() string
	GetToken() string
	GetVersion() string
	GetConnectionConfig() ssh_handler.SshConfig

	IsValid() error
}

func isNodeConfigValid(nodeConfig NodeConfigInterface) error {
	var err error

	config := nodeConfig.GetConnectionConfig()
	err = config.IsValid()
	if err != nil {
		return err
	}

	if nodeConfig.GetHost() == "" {
		return fmt.Errorf("host is empty")
	}

	if nodeConfig.GetToken() == "" {
		return fmt.Errorf("token is empty")
	}

	if nodeConfig.GetVersion() == "" {
		return fmt.Errorf("version is empty")
	}

	return nil
}

type K3sMasterNodeConfig struct {
	host    string
	token   string
	version string

	connectionConfig *ssh_handler.SshConfig
}

func NewK3sMasterNodeConfig(host string, token string, version string, connectionConfig *ssh_handler.SshConfig) *K3sMasterNodeConfig {
	return &K3sMasterNodeConfig{
		host:             host,
		token:            token,
		version:          version,
		connectionConfig: connectionConfig,
	}
}

func (k K3sMasterNodeConfig) GetHost() string {
	return k.host
}
func (k K3sMasterNodeConfig) GetToken() string {
	return k.token
}
func (k K3sMasterNodeConfig) GetVersion() string {
	return k.version
}
func (k K3sMasterNodeConfig) GetConnectionConfig() ssh_handler.SshConfig {
	return *k.connectionConfig
}
func (k K3sMasterNodeConfig) IsValid() error {
	return isNodeConfigValid(k)
}

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
func (k K3sWorkerNodeConfig) GetConnectionConfig() ssh_handler.SshConfig {
	return *k.connectionConfig
}
func (k K3sWorkerNodeConfig) IsValid() error {
	if k.GetServer() == "" {
		return fmt.Errorf("server is empty")
	}
	return isNodeConfigValid(k)
}
