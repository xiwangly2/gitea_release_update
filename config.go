package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	URL           string   `yaml:"url"`
	Token         string   `yaml:"token"`
	DownloadDir   string   `yaml:"download_dir"`
	DownloadFiles []string `yaml:"download_files"`
}

func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
