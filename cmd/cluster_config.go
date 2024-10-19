package cmd

type NodeConfig struct {
	Name       string `yaml:"name"`
	Connection struct {
		Host                 string `yaml:"host"`
		Port                 string `yaml:"port"`
		User                 string `yaml:"user"`
		Password             string `yaml:"password"`
		PrivateKey           string `yaml:"private_key"`
		PrivateKeyPassphrase string `yaml:"private_key_passphrase"`
	} `yaml:"connection"`
	Options []string `yaml:"options"`
}

type CusterConfig struct {
	Cluster struct {
		Version       string `yaml:"version"`
		Token         string `yaml:"token"`
		ServerAddress string `yaml:"server_address"`
	} `yaml:"cluster"`

	MasterNodes []NodeConfig `yaml:"master_nodes"`

	WorkerNodes []NodeConfig `yaml:"worker_nodes"`
}
