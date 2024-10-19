package ssh_handler

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"log"
	"testing"
)

func TestNewSShHandlerFromPassword(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	_, err := NewSshHandler(
		NewSshConfig(host, port, user, password, "", ""),
	)
	if err != nil {
		t.Errorf("Error creating SSH handler: %s", err)
		return
	}
}

func TestNewSshHandlerFromPrivateKey(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	var passphrase string

	sshHandler, err := NewSshHandler(
		NewSshConfig(host, port, user, password, "", ""),
	)

	privateKeyBytes, err := copyPublicKeyToServer(sshHandler, passphrase)
	if err != nil {
		t.Errorf("Error copying public key to server: %s", err)
		return
	}

	privateKey := string(*privateKeyBytes)

	_, err = NewSshHandler(
		NewSshConfig(host, port, user, password, privateKey, passphrase),
	)
	if err != nil {
		t.Errorf("Error creating SSH handler: %s", err)
		return
	}
}

func TestNewSshHandlerFromPrivateKeyWithPassphrase(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"
	passphrase := "test_passphrase"

	sshHandler, err := NewSshHandler(
		NewSshConfig(host, port, user, password, "", ""),
	)

	privateKeyBytes, err := copyPublicKeyToServer(sshHandler, passphrase)
	if err != nil {
		t.Errorf("Error copying public key to server: %s", err)
		return
	}

	privateKey := string(*privateKeyBytes)

	_, err = NewSshHandler(
		NewSshConfig(host, port, user, password, privateKey, passphrase),
	)
	if err != nil {
		t.Errorf("Error creating SSH handler: %s", err)
		return
	}
}

func TestSshHandler_WithSessionReturning(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	sshHandler, err := NewSshHandler(
		NewSshConfig(host, port, user, password, "", ""),
	)
	if err != nil {
		t.Errorf("Error creating SSH handler: %s", err)
		return
	}

	command := &SshCommand{
		BaseCommand: "echo",
		Args:        []string{"'hello world'"},
	}
	output, err := sshHandler.WithSessionReturning(command, &bytes.Buffer{})
	if err != nil && string(output[:]) != "hello world" {
		t.Errorf("Error running command: %s", err)
	}
}

func TestSSHHandler_WithSession(t *testing.T) {
	host := "localhost"
	port := "2222"
	user := "sshuser"
	password := "password"

	sshHandler, err := NewSshHandler(
		NewSshConfig(host, port, user, password, "", ""),
	)
	if err != nil {
		t.Errorf("Error creating SSH handler: %s", err)
		return
	}

	command := &SshCommand{
		BaseCommand: "echo",
		Args:        []string{"'hello world'"},
	}
	err = sshHandler.WithSession(command, &bytes.Buffer{})
	if err != nil {
		t.Errorf("Error running command: %s", err)
	}
}

// Auxiliar functions

func copyPublicKeyToServer(sshHandler *SshHandler, passphrase string) (*[]byte, error) {
	privateKey, err := generatePrivateKey(4096)
	if err != nil {
		return nil, err
	}

	publicKey, err := generatePublicKey(&privateKey.PublicKey)
	pemBytes, err := encodePrivateKeyToPEM(privateKey, passphrase)

	err = sshHandler.WithSession(&SshCommand{
		BaseCommand: "mkdir -p ~/.ssh;",
	}, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	err = sshHandler.WithSession(&SshCommand{
		BaseCommand: "cat >> ~/.ssh/authorized_keys",
	}, bytes.NewBuffer(publicKey))
	if err != nil {
		return nil, err
	}

	return &pemBytes, nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey, pwd string) ([]byte, error) {
	var err error
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Encrypt the pem
	if pwd != "" {
		block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(pwd), x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}
	}

	return pem.EncodeToMemory(block), nil
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}
