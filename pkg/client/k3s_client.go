package client

import (
	"bytes"
	"errors"
	"fmt"
	"yoshi_k3s/pkg/resources"
	"yoshi_k3s/pkg/ssh_handler"
)

type K3sClient struct {
	k3sCommandPrefix string
	k3sBaseCommand   string

	masterNodes []resources.K3sMasterNodeConfig
	workerNodes []resources.K3sWorkerNodeConfig
}

func NewK3sClient() *K3sClient {
	return &K3sClient{
		k3sCommandPrefix: "curl -sfL https://get.k3s.io | ",
		k3sBaseCommand:   "sh -",
	}
}

func (c *K3sClient) ConfigureMasterNode(k3sConfig resources.K3sMasterNodeConfig, options []string) error {
	err := c.validateNodeConfig(k3sConfig)
	if err != nil {
		return err
	}

	var envVariablesMap = make(map[string]string)

	envVariablesMap["INSTALL_K3S_EXEC"] = "server"

	return c.configureNode(k3sConfig, envVariablesMap, options)
}

func (c *K3sClient) ConfigureWorkerNode(k3sConfig resources.K3sWorkerNodeConfig, options []string) error {
	err := c.validateWorkerNodeConfig(k3sConfig)
	if err != nil {
		return err
	}

	var envVariablesMap = make(map[string]string)

	envVariablesMap["INSTALL_K3S_EXEC"] = "agent"
	envVariablesMap["K3S_URL"] = fmt.Sprintf("https://%s:6443", k3sConfig.Server)

	return c.configureNode(k3sConfig.K3sMasterNodeConfig, envVariablesMap, options)
}

func (c *K3sClient) configureNode(k3sConfig resources.K3sMasterNodeConfig,
	envVariablesMap map[string]string,
	options []string) error {
	envVariablesMap["INSTALL_K3S_VERSION"] = k3sConfig.Version
	envVariablesMap["K3S_TOKEN"] = k3sConfig.Token

	var sshCommand = ssh_handler.SshCommand{
		CommandPrefix: c.k3sCommandPrefix,
		BaseCommand:   c.k3sBaseCommand,
		EnvVars:       envVariablesMap,
		Args:          options,
	}

	sshHandler, err := c.createSshHandler(k3sConfig.ConnectionConfig)
	if err != nil {
		return err
	}

	output, err := sshHandler.WithSession(&sshCommand, bytes.Buffer{})
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil
}

func (c *K3sClient) createSshHandler(sshConfig ssh_handler.SshConfig) (*ssh_handler.SSHHandler, error) {
	if sshConfig.Password != "" {
		return ssh_handler.NewSShHandlerFromPassword(sshConfig.Host, sshConfig.Port, sshConfig.User, sshConfig.Password)
	} else if sshConfig.PrivateKeyPassphrase != "" {
		return ssh_handler.NewSshHandlerFromPrivateKeyWithPassphrase(sshConfig.Host, sshConfig.Port, sshConfig.User,
			sshConfig.PrivateKey, sshConfig.PrivateKeyPassphrase)
	} else {
		return ssh_handler.NewSShHandlerFromPrivateKey(sshConfig.Host, sshConfig.Port, sshConfig.User, sshConfig.PrivateKey)
	}
}

func (c *K3sClient) validateNodeConfig(nodeConfig resources.K3sMasterNodeConfig) error {
	if nodeConfig.Host == "" {
		return errors.New("host is empty")
	}

	if nodeConfig.Token == "" {
		return errors.New("token is empty")
	}

	if nodeConfig.Version == "" {
		return errors.New("version is empty")
	}

	return c.validateNodeConnection(nodeConfig.ConnectionConfig)
}

func (c *K3sClient) validateWorkerNodeConfig(nodeConfig resources.K3sWorkerNodeConfig) error {
	if nodeConfig.Server == "" {
		return errors.New("server is empty")
	}

	return c.validateNodeConfig(nodeConfig.K3sMasterNodeConfig)
}

func (c *K3sClient) validateNodeConnection(nodeConnection ssh_handler.SshConfig) error {
	if nodeConnection.Host == "" {
		return errors.New("host is empty")
	}

	if nodeConnection.User == "" {
		return errors.New("user is empty")
	}

	if nodeConnection.PrivateKey == "" && nodeConnection.Password == "" {
		return errors.New("either privateKey or password must be set")
	}

	return nil
}
