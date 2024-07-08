package ssh_handler

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

type SshOptions struct {
	WithPTY bool
}

type SSHHandler struct {
	sshClient *ssh.Client

	sshContext context.Context
}

func NewSshHandler(sshConfig *SshConfig) (*SSHHandler, error) {
	sshClient, err := createNewSshClient(sshConfig)
	if err != nil {
		return nil, err
	}

	return &SSHHandler{
		sshClient: sshClient,
	}, nil
}

func (s *SSHHandler) WithOptions(options *SshOptions) (*context.Context, *context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	s.sshContext = context.WithValue(ctx, "options", options)

	return &ctx, &cancelFunc
}

func (s *SSHHandler) WithSession(ssCommand SshCommandInterface, input *bytes.Buffer) error {
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

func (s *SSHHandler) createSshSession() (*ssh.Session, error) {
	session, err := s.sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	var options *SshOptions
	if s.sshContext != nil {
		options = s.sshContext.Value("options").(*SshOptions)
	}

	if options == nil || options.WithPTY {
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
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
