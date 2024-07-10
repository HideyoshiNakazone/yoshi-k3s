package resources

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"testing"
)

func TestK3sWorkerNodeConfig_IsValid(t *testing.T) {
	n := K3sWorkerNodeConfig{
		server: "master_node",

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
