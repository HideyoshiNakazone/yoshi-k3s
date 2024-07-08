package ssh_handler

import "fmt"

type SshConfig struct {
	host                 string
	port                 string
	user                 string
	password             string
	privateKey           string
	privateKeyPassphrase string
}

func NewSshConfig(host string, port string, user string, password string, privateKey string, privateKeyPassphrase string) *SshConfig {
	return &SshConfig{
		host:                 host,
		port:                 port,
		user:                 user,
		password:             password,
		privateKey:           privateKey,
		privateKeyPassphrase: privateKeyPassphrase,
	}
}

func (s *SshConfig) GetHost() string {
	return s.host
}

func (s *SshConfig) GetPort() string {
	return s.port
}

func (s *SshConfig) GetUser() string {
	return s.user
}

func (s *SshConfig) GetPassword() string {
	return s.password
}

func (s *SshConfig) GetPrivateKey() string {
	return s.privateKey
}

func (s *SshConfig) GetPrivateKeyPassphrase() string {
	return s.privateKeyPassphrase
}

func (s *SshConfig) IsValid() error {
	if s.host == "" {
		return fmt.Errorf("host is empty")
	}

	if s.port == "" {
		return fmt.Errorf("port is empty")
	}

	if s.user == "" {
		return fmt.Errorf("user is empty")
	}

	if s.password == "" && s.privateKey == "" {
		return fmt.Errorf("password and private key are empty")
	}

	return nil
}
