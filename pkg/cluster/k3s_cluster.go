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

	k3sVersion       string
	k3sKubeConfig    string
	k3sToken         string
	k3sServerAddress string
}

func NewK3sClient(token string, serverAddress string) *K3sCluster {
	if token == "" || serverAddress == "" {
		return nil
	}

	return &K3sCluster{
		k3sCommandPrefix: "curl -sfL https://get.k3s.io |",
		k3sBaseCommand:   "sh -s -",
		k3sToken:         token,
		k3sServerAddress: serverAddress,
	}
}

func NewK3sClientWithVersion(version string, token string, serverAddress string) *K3sCluster {
	client := NewK3sClient(token, serverAddress)
	client.k3sVersion = version

	return client
}

func (c *K3sCluster) ConfigureMasterNode(k3sConfig resources.NodeConfig, options []string) (*[]byte, error) {
	err := k3sConfig.IsValid()
	if err != nil {
		return nil, err
	}

	options = append([]string{"server"}, options...)

	var envVariablesMap = make(map[string]string)
	envVariablesMap["K3S_KUBECONFIG_MODE"] = "644"

	err = c.configureNode(k3sConfig, envVariablesMap, options)
	if err != nil {
		return nil, err
	}

	kubeconfig, err := c.configureKubeconfig(k3sConfig.GetConnectionConfig())
	if err != nil {
		return nil, err
	}

	return &kubeconfig, err
}

func (c *K3sCluster) ConfigureWorkerNode(k3sConfig resources.NodeConfig, options []string) error {
	err := k3sConfig.IsValid()
	if err != nil {
		return err
	}

	var envVariablesMap = make(map[string]string)
	envVariablesMap["K3S_URL"] = fmt.Sprintf("https://%s:6443", c.k3sServerAddress)

	options = append([]string{"agent"}, options...)

	return c.configureNode(k3sConfig, envVariablesMap, options)
}

func (c *K3sCluster) DestroyMasterNode(k3sConfig resources.NodeConfig) error {
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

func (c *K3sCluster) DestroyWorkerNode(k3sConfig resources.NodeConfig) error {
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

func (c *K3sCluster) configureNode(k3sConfig resources.NodeConfig,
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

func (c *K3sCluster) configureKubeconfig(connectionConfig *ssh_handler.SshConfig) ([]byte, error) {
	sshHandler, err := ssh_handler.NewSshHandler(connectionConfig)
	if err != nil {
		return []byte(""), err
	}

	commands := []string{
		"mkdir -p $HOME/.kube;",
		"cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config;",
		"chmod g+r $HOME/.kube/config;",
	}

	err = c.executeK3sCommand(
		sshHandler,
		&ssh_handler.SshCommand{
			BaseCommand: strings.Join(commands, " "),
		},
		connectionConfig.GetPassword(),
	)
	if err != nil {
		return []byte(""), err
	}

	kubeconfigContent, err := c.executeK3sCommandWithOutput(
		sshHandler,
		&ssh_handler.SshCommand{
			BaseCommand: "cat $HOME/.kube/config",
		},
		connectionConfig.GetPassword(),
	)

	if err != nil {
		return []byte(""), err
	}

	LOCALHOST_ADDRESS := "127.0.0.1"
	kubeconfigContent = bytes.Replace(
		kubeconfigContent,
		[]byte(LOCALHOST_ADDRESS),
		[]byte(c.k3sServerAddress),
		1,
	)

	return kubeconfigContent, err
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

func (c *K3sCluster) executeK3sCommandWithOutput(sshHandler *ssh_handler.SshHandler,
	command *ssh_handler.SshCommand,
	password string) ([]byte, error) {
	terminalMode := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if ctxt, cancelFunc := sshHandler.WithTerminalMode(&terminalMode); ctxt != nil {
		defer (*cancelFunc)()
		return sshHandler.WithSessionReturning(
			command,
			bytes.NewBuffer([]byte(password+"\n")),
		)
	}

	return nil, errors.New("failed to configure ssh session")
}
