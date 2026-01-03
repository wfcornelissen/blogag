package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func getHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir, nil
}

func (c *Config) SetUser(user string) error {
	newCfg := Config{
		DbUrl:           c.DbUrl,
		CurrentUserName: user,
	}

	byteData, err := json.Marshal(newCfg)
	if err != nil {
		return err
	}

	homeDir, err := getHomeDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(homeDir+"/"+configFilePath, byteData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory")
		return Config{}, err
	}

	data, err := os.ReadFile(homeDir + "/" + configFilePath)
	if err != nil {
		fmt.Println("Error reading from file")
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
