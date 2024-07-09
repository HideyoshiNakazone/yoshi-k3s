package resources

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sNode_Init_Checks_If_Valid(t *testing.T) {
	n := K3sMasterNodeConfig{
		host:    "127.0.0.1",
		token:   "token",
		version: "latest",
		connectionConfig: ssh_handler.NewSshConfig(
			"host",
			"port",
			"user",
			"password",
			"",
			"",
		),
	}

	if err := isNodeConfigValid(n); err != nil {
		t.Errorf("Expected valid node config, got error: %s", err)
	}
}

func TestK3sMasterNodeConfig_HasChanged(t *testing.T) {
	n := K3sMasterNodeConfig{
		host:    "127.0.0.1",
		token:   "token",
		version: "latest",
		connectionConfig: ssh_handler.NewSshConfig(
			"host",
			"port",
			"user",
			"password",
			"",
			"",
		),
	}

	other := K3sMasterNodeConfig{
		host:    "127.0.0.1",
		token:   "token",
		version: "latest",
		connectionConfig: ssh_handler.NewSshConfig(
			"host",
			"port",
			"user",
			"password",
			"",
			"",
		),
	}

	if n.HasChanged(other) {
		t.Errorf("Expected no change, got change")
	}

	other.token = "new_token"

	if !n.HasChanged(other) {
		t.Errorf("Expected change, got no change")
	}
}
