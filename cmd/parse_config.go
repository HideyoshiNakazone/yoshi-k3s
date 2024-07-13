package cmd

import (
	"fmt"
	"os"
)

type CusterConfig struct {
	version string
}

func ParseConfig(configPath string) *CusterConfig {
	content, err := readFileContent(configPath)
	if err != nil {
		return nil
	}

	fmt.Println(content)

	return nil
}

func readFileContent(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
