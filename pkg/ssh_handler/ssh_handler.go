package ssh_handler

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

type SshHandler struct {
	sshClient *ssh.Client

	sshContext     context.Context
	sshContextFlag string
}

func NewSshHandler(sshConfig *SshConfig) (*SshHandler, error) {
	sshClient, err := createNewSshClient(sshConfig)
	if err != nil {
		return nil, err
	}

	return &SshHandler{
		sshClient:      sshClient,
		sshContextFlag: "terminalModes",
	}, nil
}

func (s *SshHandler) WithTerminalMode(modes *ssh.TerminalModes) (*context.Context, *context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	s.sshContext = context.WithValue(ctx, s.sshContextFlag, modes)

	return &ctx, &cancelFunc
}

func (s *SshHandler) WithSession(ssCommand SshCommandInterface, input *bytes.Buffer) error {
	session, err := s.createSshSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if input.Len() > 0 {
		session.Stdin = input
	}

	command, err := ssCommand.GetParsedCommand()
	if err != nil {
		return err
	}

	return session.Run(command)
}

func (s *SshHandler) WithSessionReturning(ssCommand SshCommandInterface, input *bytes.Buffer) ([]byte, error) {
	session, err := s.createSshSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	if input.Len() > 0 {
		session.Stdin = input
	}

	command, err := ssCommand.GetParsedCommand()
	if err != nil {
		return nil, err
	}

	return session.Output(command)
}

func (s *SshHandler) Close() error {
	return s.sshClient.Close()
}

func (s *SshHandler) createSshSession() (*ssh.Session, error) {
	session, err := s.sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	var terminalModes *ssh.TerminalModes
	if s.sshContext != nil {
		terminalModes = s.sshContext.Value(s.sshContextFlag).(*ssh.TerminalModes)
	}

	if terminalModes != nil {
		if err := session.RequestPty("xterm", 80, 40, *terminalModes); err != nil {
			session.Close()
			return nil, err
		}
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session, nil
}

func createNewSshClient(sshConfig *SshConfig) (*ssh.Client, error) {
	if sshConfig.GetPassword() != "" {
		return newSShClientFromPassword(sshConfig)
	} else if sshConfig.GetPrivateKeyPassphrase() != "" {
		return newSshClientFromPrivateKeyWithPassphrase(sshConfig)
	} else {
		return newSShClientFromPrivateKey(sshConfig)
	}
}

func newSShClientFromPrivateKey(sshConfig *SshConfig) (*ssh.Client, error) {
	hostAddrString, err := parseAddrString(sshConfig.GetHost(), sshConfig.GetPort())
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey([]byte(sshConfig.GetPrivateKey()))
	if err != nil {
		return nil, err
	}

	return ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            sshConfig.GetUser(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	})
}

func newSshClientFromPrivateKeyWithPassphrase(sshConfig *SshConfig) (*ssh.Client, error) {
	hostAddrString, err := parseAddrString(sshConfig.GetHost(), sshConfig.GetPort())
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(
		[]byte(sshConfig.GetPrivateKey()), []byte(sshConfig.GetPrivateKeyPassphrase()),
	)
	if err != nil {
		return nil, err
	}

	return ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            sshConfig.GetUser(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	})
}

func newSShClientFromPassword(sshConfig *SshConfig) (*ssh.Client, error) {
	hostAddrString, err := parseAddrString(sshConfig.GetHost(), sshConfig.GetPort())
	if err != nil {
		return nil, err
	}

	return ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            sshConfig.GetUser(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConfig.GetPassword()),
		},
	})
}

func parseAddrString(host string, port string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("host is empty")
	}

	if port == "" {
		return "", fmt.Errorf("port is empty")
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
