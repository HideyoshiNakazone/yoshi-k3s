package resources

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"fmt"
)

type NodeConfigInterface interface {
	GetHost() string
	GetToken() string
	GetVersion() string
	GetConnectionConfig() *ssh_handler.SshConfig

	IsValid() error
	HasChanged(other NodeConfigInterface) bool
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

func hasChanged(nodeConfig NodeConfigInterface, other NodeConfigInterface) bool {
	if nodeConfig.GetHost() != other.GetHost() {
		return true
	}

	if nodeConfig.GetToken() != other.GetToken() {
		return true
	}

	if nodeConfig.GetVersion() != other.GetVersion() {
		return true
	}

	if nodeConfig.GetConnectionConfig().HasChanged(other.GetConnectionConfig()) {
		return true
	}

	return false
}
