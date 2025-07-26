package main

import (
	"encoding/json"
	"os"
	"path"
)

func (m *App) LoadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := path.Join(configDir, "phoningdl", "config.json")
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
	configPath := path.Join(configDir, "phoningdl", "config.json")
	os.MkdirAll(path.Join(configDir, "phoningdl"), os.ModePerm)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m.config)
}