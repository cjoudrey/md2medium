package main

import "os"
import "os/user"
import "path"
import "encoding/json"
import "io/ioutil"

type Config struct {
	GitHubAccessToken string `json:"github_access_token"`
	MediumAccessToken string `json:"medium_access_token"`
	MediumUserId      string `json:"medium_user_id"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	configPath, err := configPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return nil, err
	}

	configContents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(configContents, config)

	return config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := configPath()
	if err != nil {
		return err
	}

	configJson, _ := json.Marshal(config)

	if err := ioutil.WriteFile(configPath, configJson, 0644); err != nil {
		return err
	}

	return nil
}

func configPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(user.HomeDir, ".md2medium"), nil
}
