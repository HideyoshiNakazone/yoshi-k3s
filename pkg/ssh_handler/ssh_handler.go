package ssh_handler

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type SSHHandler struct {
	sshClient *ssh.Client
}

func (s *SSHHandler) WithSession(ssCommand SshCommandInterface, input bytes.Buffer) (string, error) {
	var session *ssh.Session
	var output bytes.Buffer
	var err error

	session, err = s.sshClient.NewSession()
	if err != nil {
		return output.String(), err
	}
	defer session.Close()

	if input.Len() > 0 {
		session.Stdin = &input
	}
	session.Stdout = &output

	var command string
	command, err = ssCommand.GetParsedCommand()
	if err != nil {
		return "", err
	}

	err = session.Run(command)
	return output.String(), err
}

func NewSShHandlerFromPrivateKey(host string, port string, user string, privateKey string) (*SSHHandler, error) {
	hostAddrString, err := parseAddrString(host, port)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SSHHandler{sshClient: sshClient}, nil
}

func NewSshHandlerFromPrivateKeyWithPassphrase(host string, port string, user string, privateKey string, passphrase string) (*SSHHandler, error) {
	hostAddrString, err := parseAddrString(host, port)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passphrase))
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SSHHandler{sshClient: sshClient}, nil
}

func NewSShHandlerFromPassword(host string, port string, user string, password string) (*SSHHandler, error) {
	hostAddrString, err := parseAddrString(host, port)
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.Dial("tcp", hostAddrString, &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	})
	if err != nil {
		return nil, err
	}

	return &SSHHandler{sshClient: sshClient}, nil
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
