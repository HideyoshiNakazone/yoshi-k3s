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

	err := isNodeConfigValid(n)
	if err != nil {
		t.Errorf("Expected valid node config, got error: %s", err)
	}
}
