package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabaseUrl string `json:"db_url"`
	Username    string `json:"current_user_name"`
}

func ReadConfig() Config {

	path, _ := getConfigFilePath()

	file, err := os.Open(path)
	if err != nil {
		return Config{}
	}
	defer file.Close()

	config := Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}
	}

	return config
}

func (config *Config) SetUser(user string) error {
	config.Username = user
	write(*config)
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := home + "/" + configFileName

	return path, nil
}

func write(config Config) error {
	path, _ := getConfigFilePath()

	file, _ := os.Create(path)
	defer file.Close()

	configData, _ := json.MarshalIndent(config, "", " ")

	file.Write(configData)

	return nil
}
