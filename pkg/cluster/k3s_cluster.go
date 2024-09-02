package cluster

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/resources"
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"golang.org/x/crypto/ssh"
	"strings"
)

type K3sCluster struct {
	k3sCommandPrefix string
	k3sBaseCommand   string

	k3sVersion string
	k3sToken   string
}

func NewK3sClient(token string) *K3sCluster {
	if token == "" {
		return nil
	}

	return &K3sCluster{
		k3sCommandPrefix: "curl -sfL https://get.k3s.io |",
		k3sBaseCommand:   "sh -s -",
		k3sToken:         token,
	}
}

func NewK3sClientWithVersion(version string, token string) *K3sCluster {
	client := NewK3sClient(token)
	client.k3sVersion = version

	return client
}

func (c *K3sCluster) ConfigureMasterNode(k3sConfig resources.K3sMasterNodeConfig, options []string) error {
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
		"cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config;",
		"chmod g+r $HOME/.kube/config;",
	}

	config := k3sConfig.GetConnectionConfig()
	err = c.executeK3sCommand(
		sshHandler,
		&ssh_handler.SshCommand{
			BaseCommand: strings.Join(commands, " "),
		},
		config.GetPassword(),
	)

	return err
}

func (c *K3sCluster) ConfigureWorkerNode(k3sConfig resources.K3sWorkerNodeConfig, options []string) error {
	err := k3sConfig.IsValid()
	if err != nil {
		return err
	}

	var envVariablesMap = make(map[string]string)
	envVariablesMap["K3S_URL"] = fmt.Sprintf("https://%s:6443", k3sConfig.GetServer())

	options = append([]string{"agent"}, options...)

	return c.configureNode(k3sConfig, envVariablesMap, options)
}

func (c *K3sCluster) DestroyMasterNode(k3sConfig resources.K3sMasterNodeConfig) error {
	sshHandler, err := ssh_handler.NewSshHandler(k3sConfig.GetConnectionConfig())
	if err != nil {
		return err
	}

	config := k3sConfig.GetConnectionConfig()
	err = c.executeK3sCommand(
		sshHandler,
		&ssh_handler.SshCommand{
			BaseCommand: "sudo k3s-uninstall.sh",
		},
		config.GetPassword(),
	)

	return err
}

func (c *K3sCluster) DestroyWorkerNode(k3sConfig resources.K3sWorkerNodeConfig) error {
	sshHandler, err := ssh_handler.NewSshHandler(k3sConfig.GetConnectionConfig())
	if err != nil {
		return err
	}

	config := k3sConfig.GetConnectionConfig()
	err = c.executeK3sCommand(
		sshHandler,
		&ssh_handler.SshCommand{
			BaseCommand: "sudo k3s-agent-uninstall.sh",
		},
		config.GetPassword(),
	)

	return err
}

func (c *K3sCluster) configureNode(k3sConfig resources.NodeConfigInterface,
	envVariablesMap map[string]string,
	options []string) error {
	envVariablesMap["K3S_TOKEN"] = c.k3sToken

	if c.k3sVersion != "" {
		envVariablesMap["INSTALL_K3S_VERSION"] = c.k3sVersion
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
	return c.executeK3sCommand(
		sshHandler,
		&sshCommandCreateNode,
		config.GetPassword(),
	)
}

func (c *K3sCluster) executeK3sCommand(sshHandler *ssh_handler.SshHandler,
	command *ssh_handler.SshCommand,
	password string) error {
	terminalMode := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if ctxt, cancelFunc := sshHandler.WithTerminalMode(&terminalMode); ctxt != nil {
		defer (*cancelFunc)()
		return sshHandler.WithSession(
			command,
			bytes.NewBuffer([]byte(password+"\n")),
		)
	}

	return errors.New("failed to configure ssh session")
}
