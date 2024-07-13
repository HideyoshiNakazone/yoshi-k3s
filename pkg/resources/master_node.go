package resources

import "github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"

type K3sMasterNodeConfig struct {
	connectionConfig *ssh_handler.SshConfig
}

// Check if K3sMasterNodeConfig implements NodeConfigInterface
var _ NodeConfigInterface = &K3sMasterNodeConfig{}

func NewK3sMasterNodeConfig(connectionConfig *ssh_handler.SshConfig) *K3sMasterNodeConfig {
	return &K3sMasterNodeConfig{
		connectionConfig: connectionConfig,
	}
}

func (k K3sMasterNodeConfig) GetConnectionConfig() *ssh_handler.SshConfig {
	return k.connectionConfig
}
func (k K3sMasterNodeConfig) IsValid() error {
	return isNodeConfigValid(k)
}
