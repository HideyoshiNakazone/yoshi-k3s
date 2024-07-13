package main

import (
	"flag"
	"github.com/HideyoshiNakazone/yoshi-k3s/cmd"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "str", "config.yml", "Path to Config File [default=config.yml]")
	flag.Parse()

	cmd.ParseConfig(configPath)
}
