package resources

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sNode_Init_Checks_If_Valid(t *testing.T) {
	n := K3sMasterNodeConfig{
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
