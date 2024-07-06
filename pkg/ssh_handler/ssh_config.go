package ssh_handler

type SshConfig struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	PrivateKey           string
	PrivateKeyPassphrase string
}
