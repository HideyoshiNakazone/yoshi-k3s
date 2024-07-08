package resources

import "terraform-yoshi-k3s/pkg/ssh_handler"

type NodeConfigInterface interface {
	GetHost() string
	GetToken() string
	GetVersion() string
	GetConnectionConfig() ssh_handler.SshConfig
}

type K3sMasterNodeConfig struct {
	Host    string
	Token   string
	Version string

	ConnectionConfig ssh_handler.SshConfig
}

func (k K3sMasterNodeConfig) GetHost() string {
	return k.Host
}
func (k K3sMasterNodeConfig) GetToken() string {
	return k.Token
}
func (k K3sMasterNodeConfig) GetVersion() string {
	return k.Version
}
func (k K3sMasterNodeConfig) GetConnectionConfig() ssh_handler.SshConfig {
	return k.ConnectionConfig
}

type K3sWorkerNodeConfig struct {
	Server string

	Host    string
	Token   string
	Version string

	ConnectionConfig ssh_handler.SshConfig
}

func (k K3sWorkerNodeConfig) GetServer() string {
	return k.Server
}
func (k K3sWorkerNodeConfig) GetHost() string {
	return k.Host
}
func (k K3sWorkerNodeConfig) GetToken() string {
	return k.Token
}
func (k K3sWorkerNodeConfig) GetVersion() string {
	return k.Version
}
func (k K3sWorkerNodeConfig) GetConnectionConfig() ssh_handler.SshConfig {
	return k.ConnectionConfig
}
