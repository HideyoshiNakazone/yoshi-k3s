package client

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/resources"
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type NodeMapping[NodeType resources.NodeConfigInterface] map[string]NodeType

type K3sClient struct {
	k3sCommandPrefix string
	k3sBaseCommand   string

	masterNodes NodeMapping[resources.K3sMasterNodeConfig]
	workerNodes NodeMapping[resources.K3sWorkerNodeConfig]
}

func NewK3sClient(masterNodes NodeMapping[resources.K3sMasterNodeConfig],
	workerNodes NodeMapping[resources.K3sWorkerNodeConfig]) *K3sClient {
	return &K3sClient{
		k3sCommandPrefix: "curl -sfL https://get.k3s.io |",
		k3sBaseCommand:   "sh -s -",
		masterNodes:      masterNodes,
		workerNodes:      workerNodes,
	}
}

func (c *K3sClient) ConfigureMasterNode(k3sConfig resources.K3sMasterNodeConfig, options []string) error {
	err := k3sConfig.IsValid()
	if err != nil {
		return err
	}

	options = append([]string{"server"}, options...)

	err = c.configureNode(k3sConfig, make(map[string]string), options)
	if err != nil {
		return err
	}

	sshHandler, err := ssh_handler.NewSshHandler(k3sConfig.GetConnectionConfig())
	if err != nil {
		return err
	}

	commands := []string{
		"sudo chmod 644 /etc/rancher/k3s/k3s.yaml;",
		"mkdir -p $HOME/.kube;",
		"cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/k3s.yaml;",
		"chmod g+r $HOME/.kube/k3s.yaml;",
	}

	config := k3sConfig.GetConnectionConfig()
	err = sshHandler.WithSession(
		&ssh_handler.SshCommand{
			BaseCommand: strings.Join(commands, " "),
		},
		bytes.NewBuffer([]byte(config.GetPassword()+"\n")),
	)

	if err == nil {
		c.masterNodes[k3sConfig.GetHost()] = k3sConfig
	}

	return err
}

func (c *K3sClient) ConfigureWorkerNode(k3sConfig resources.K3sWorkerNodeConfig, options []string) error {
	if len(c.masterNodes) == 0 {
		return errors.New("no master nodes configured")
	}

	err := k3sConfig.IsValid()
	if err != nil {
		return err
	}

	var envVariablesMap = make(map[string]string)
	envVariablesMap["K3S_URL"] = fmt.Sprintf("https://%s:6443", k3sConfig.GetServer())

	options = append([]string{"agent"}, options...)

	err = c.configureNode(k3sConfig, envVariablesMap, options)
	if err == nil {
		c.workerNodes[k3sConfig.GetHost()] = k3sConfig
	}

	return err
}

func (c *K3sClient) configureNode(k3sConfig resources.NodeConfigInterface,
	envVariablesMap map[string]string,
	options []string) error {
	envVariablesMap["K3S_TOKEN"] = k3sConfig.GetToken()

	if k3sConfig.GetVersion() != "" {
		envVariablesMap["INSTALL_K3S_VERSION"] = k3sConfig.GetVersion()
	}

	var sshCommandCreateNode = ssh_handler.SshCommand{
		CommandPrefix: c.k3sCommandPrefix,
		BaseCommand:   c.k3sBaseCommand,
		EnvVars:       envVariablesMap,
		Args:          options,
	}

	sshHandler, err := ssh_handler.NewSshHandler(k3sConfig.GetConnectionConfig())
	if err != nil {
		return err
	}

	config := k3sConfig.GetConnectionConfig()
	err = sshHandler.WithSession(
		&sshCommandCreateNode,
		bytes.NewBuffer([]byte(config.GetPassword()+"\n")),
	)
	return err
}
