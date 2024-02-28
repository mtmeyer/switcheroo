package utils

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
)

type Config struct {
	Directories map[string]string `json:"directories"`
	Output      string            `json:"output"`
	Plugins     []string          `json:"plugins"`
}

var ConfigDirectory string

func GetConfigFileDirectory() (string, error) {
	const CONFIG_DIR = "/switcheroo"

	var xdgHome = os.Getenv("XDG_CONFIG_HOME")
	if len(xdgHome) > 0 {
		return path.Join(xdgHome, CONFIG_DIR), nil
	}

	var home = os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, "/.config", CONFIG_DIR), nil
	}

	return "", errors.New("No config directory found")
}

func ParseConfig() (*Config, error) {
	var configFileErr error
	ConfigDirectory, configFileErr = GetConfigFileDirectory()

	configPath := path.Join(ConfigDirectory, "/config.json")

	if configFileErr != nil {
		return nil, configFileErr
	}

	jsonFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var config Config

	json.Unmarshal(jsonBytes, &config)

	return &config, nil
}
