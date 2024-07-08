package client

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"terraform-yoshi-k3s/pkg/resources"
	"terraform-yoshi-k3s/pkg/ssh_handler"
)

type K3sClient struct {
	k3sCommandPrefix string
	k3sBaseCommand   string

	masterNodes []resources.K3sMasterNodeConfig
	workerNodes []resources.K3sWorkerNodeConfig
}

func NewK3sClient() *K3sClient {
	return &K3sClient{
		k3sCommandPrefix: "curl -sfL https://get.k3s.io |",
		k3sBaseCommand:   "sh -s -",
	}
}

func (c *K3sClient) ConfigureMasterNode(k3sConfig resources.K3sMasterNodeConfig, options []string) error {
	err := c.validateNodeConfig(k3sConfig)
	if err != nil {
		return err
	}

	options = append([]string{"server"}, options...)

	err = c.configureNode(k3sConfig, make(map[string]string), options)
	if err != nil {
		return err
	}

	sshHandler, err := c.createSshHandler(k3sConfig.ConnectionConfig)
	if err != nil {
		return err
	}

	commands := []string{
		"sudo chmod 644 /etc/rancher/k3s/k3s.yaml;",
		"mkdir -p $HOME/.kube;",
		"cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/k3s.yaml;",
		"chmod g+r $HOME/.kube/k3s.yaml;",
	}

	output, err := sshHandler.WithSession(
		&ssh_handler.SshCommand{
			BaseCommand: strings.Join(commands, " "),
		},
		*bytes.NewBuffer([]byte(k3sConfig.ConnectionConfig.Password + "\n")),
	)
	fmt.Println(output)

	if err == nil {
		c.masterNodes = append(c.masterNodes, k3sConfig)
	}

	return err
}

func (c *K3sClient) ConfigureWorkerNode(k3sConfig resources.K3sWorkerNodeConfig, options []string) error {
	err := c.validateWorkerNodeConfig(k3sConfig)
	if err != nil {
		return err
	}

	if len(c.masterNodes) == 0 {
		return errors.New("no master nodes configured")
	}

	var envVariablesMap = make(map[string]string)
	envVariablesMap["K3S_URL"] = fmt.Sprintf("https://%s:6443", k3sConfig.Server)

	options = append([]string{"agent"}, options...)

	err = c.configureNode(k3sConfig, envVariablesMap, options)
	if err == nil {
		c.workerNodes = append(c.workerNodes, k3sConfig)
	}

	return err
}

func (c *K3sClient) configureNode(k3sConfig resources.NodeConfigInterface,
	envVariablesMap map[string]string,
	options []string) error {
	envVariablesMap["K3S_TOKEN"] = k3sConfig.GetToken()

	if k3sConfig.GetVersion() != "" {
		envVariablesMap["K3S_VERSION"] = k3sConfig.GetVersion()
	}

	var sshCommandCreateNode = ssh_handler.SshCommand{
		CommandPrefix: c.k3sCommandPrefix,
		BaseCommand:   c.k3sBaseCommand,
		EnvVars:       envVariablesMap,
		Args:          options,
	}

	sshHandler, err := c.createSshHandler(k3sConfig.GetConnectionConfig())
	if err != nil {
		return err
	}

	output, err := sshHandler.WithSession(
		&sshCommandCreateNode,
		*bytes.NewBuffer([]byte(k3sConfig.GetConnectionConfig().Password + "\n")),
	)
	fmt.Println(output)
	return err
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

func (c *K3sClient) validateNodeConfig(nodeConfig resources.NodeConfigInterface) error {
	if nodeConfig.GetHost() == "" {
		return errors.New("host is empty")
	}

	if nodeConfig.GetToken() == "" {
		return errors.New("token is empty")
	}

	return c.validateNodeConnection(nodeConfig.GetConnectionConfig())
}

func (c *K3sClient) validateWorkerNodeConfig(nodeConfig resources.K3sWorkerNodeConfig) error {
	if nodeConfig.GetServer() == "" {
		return errors.New("server is empty")
	}

	return c.validateNodeConfig(nodeConfig)
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
