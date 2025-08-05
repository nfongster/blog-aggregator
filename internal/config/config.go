package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"
const defaultDbUrl = "postgres://example"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.DbUrl = defaultDbUrl
	c.CurrentUserName = user
	if err := write(c); err != nil {
		return fmt.Errorf("error updating user to %s: %v", user, err)
	}
	return nil
}

func write(c *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := dir + string(os.PathSeparator) + configFileName
	return path, nil
}
