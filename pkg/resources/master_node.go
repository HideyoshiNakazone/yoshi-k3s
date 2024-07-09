package resources

import "github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"

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
func (k K3sMasterNodeConfig) GetConnectionConfig() *ssh_handler.SshConfig {
	return k.connectionConfig
}
func (k K3sMasterNodeConfig) IsValid() error {
	return isNodeConfigValid(k)
}
func (k K3sMasterNodeConfig) HasChanged(other NodeConfigInterface) bool {
	return hasChanged(k, other)
}
