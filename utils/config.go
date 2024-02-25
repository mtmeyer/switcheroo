package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Directories []string `json:"directories"`
}

func GetConfigFilePath() (string, error) {
	const FILE_PATH = "/switcheroo/config.json"

	var xdgHome = os.Getenv("XDG_CONFIG_HOME")
	if len(xdgHome) > 0 {
		return xdgHome + FILE_PATH, nil
	}

	var home = os.Getenv("HOME")
	if len(home) > 0 {
		return home + "/.config" + FILE_PATH, nil
	}

	return "", errors.New("No config directory found")
}

func ParseConfig() (*Config, error) {
	configPath, err := GetConfigFilePath()

	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var config Config

	json.Unmarshal(jsonBytes, &config)

	return &config, nil
}
