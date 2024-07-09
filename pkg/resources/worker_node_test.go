package resources

import (
	"HideyoshiNakazone/terraform-yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sWorkerNodeConfig_IsValid(t *testing.T) {
	n := K3sWorkerNodeConfig{
		server: "master_node",

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

	if err := n.IsValid(); err != nil {
		t.Errorf("Expected valid node config, got error: %s", err)
	}
}

func TestK3sWorkerNodeConfig_HasChanged(t *testing.T) {
	n := K3sWorkerNodeConfig{
		server: "master_node",

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

	other := K3sWorkerNodeConfig{
		server: "master_node",

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

	if n.HasChanged(&other) {
		t.Errorf("Expected no change, got change")
	}

	other.token = "new_token"
	if !n.HasChanged(&other) {
		t.Errorf("Expected change, got no change")
	}
}
