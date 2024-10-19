package main

import (
	"flag"
	"fmt"
	"github.com/HideyoshiNakazone/yoshi-k3s/cmd"
)

var (
	configPath     string
	kubeconfigPath string
)

func main() {
	flag.StringVar(&configPath, "config", "config.yml", "Path to Config File [default=config.yml]")
	flag.StringVar(&kubeconfigPath, "kubeconfig", "kubeconfig", "Path to KUBECONFIG File [default=kubeconfig]")
	flag.Bool("destroy", false, "Destroy Cluster")
	flag.Parse()

	config := cmd.ParseConfig(configPath)
	if config == nil {
		fmt.Println("Error parsing config")
		return
	}

	if flag.Lookup("destroy").Value.String() == "true" {
		err := cmd.DeleteFromConfig(config)
		if err != nil {
			fmt.Println("Error deleting cluster")
		}
	} else {
		err := cmd.ConfigureFromConfig(config, &kubeconfigPath)
		if err != nil {
			fmt.Println("Error configuring cluster")
		}
	}
}
