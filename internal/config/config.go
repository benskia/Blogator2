package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFilename string = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(cfgPath)
	defer jsonFile.Close()
	if err != nil {
		return Config{}, err
	}

	newCfg := Config{}
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&newCfg)
	if err != nil {
		return Config{}, err
	}

	return newCfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) Print() {
	fmt.Printf("DB URL: %s\n", cfg.DbURL)
	fmt.Printf("Current User Name: %s\n", cfg.CurrentUserName)
}

// helper functions

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFilename), nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(cfgPath, os.O_WRONLY, os.ModePerm)
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
