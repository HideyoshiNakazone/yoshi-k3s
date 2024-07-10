package resources

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
)

type NodeConfigInterface interface {
	GetConnectionConfig() *ssh_handler.SshConfig

	IsValid() error
}

func isNodeConfigValid(nodeConfig NodeConfigInterface) error {
	var err error

	config := nodeConfig.GetConnectionConfig()
	err = config.IsValid()
	if err != nil {
		return err
	}

	return nil
}
