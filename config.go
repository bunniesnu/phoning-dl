package main

import (
	"encoding/json"
	"os"
	"path"
)

const ConfigFileName = "config.json"

func (m *App) LoadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := path.Join(configDir, AppName, ConfigFileName)
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	config := new(Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	m.config = config
	return nil
}

func (m *App) SaveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := path.Join(configDir, AppName, ConfigFileName)
	os.MkdirAll(path.Join(configDir, AppName), os.ModePerm)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m.config)
}