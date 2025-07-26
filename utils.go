package main

import (
	"encoding/json"
	"os"
)

func LoadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configPath := configDir + "/phoningdl/config.json"
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	config := new(Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (config *Config) SaveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := configDir + "/phoningdl/config.json"
	os.MkdirAll(configDir+"/phoningdl", os.ModePerm)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}