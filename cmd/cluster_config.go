package cmd

type CusterConfig struct {
	Cluster struct {
		Version string `yaml:"version"`
		Token   string `yaml:"token"`
	} `yaml:"cluster"`

	MasterNodes []struct {
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
	} `yaml:"master_nodes"`

	WorkerNodes []struct {
		Name          string `yaml:"name"`
		ServerAddress string `yaml:"server_address"`
		Connection    struct {
			Host                 string `yaml:"host"`
			Port                 string `yaml:"port"`
			User                 string `yaml:"user"`
			Password             string `yaml:"password"`
			PrivateKey           string `yaml:"private_key"`
			PrivateKeyPassphrase string `yaml:"private_key_passphrase"`
		} `yaml:"connection"`
		Options []string `yaml:"options"`
	} `yaml:"worker_nodes"`
}
