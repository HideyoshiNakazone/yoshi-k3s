package resources

import "terraform-yoshi-k3s/pkg/ssh_handler"

type K3sMasterNodeConfig struct {
	Host    string
	Token   string
	Version string

	ConnectionConfig ssh_handler.SshConfig
}

type K3sWorkerNodeConfig struct {
	K3sMasterNodeConfig

	Server string
}
