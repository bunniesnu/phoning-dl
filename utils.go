package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
)

func (m *App) LoadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := configDir + "/phoningdl/config.json"
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
	configPath := configDir + "/phoningdl/config.json"
	os.MkdirAll(configDir+"/phoningdl", os.ModePerm)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(m.config)
}

func (m *App) UpdateState(state int) error {
	if m.w == nil {
		return fmt.Errorf("window is nil")
	}
	var screen *fyne.Container
	switch state {
	case 0: // Load config
		screen = m.LoadingConfigScreen()
	default:
		return fmt.Errorf("unknown state: %d", state)
	}
	(*m.w).SetContent(screen)
	return nil
}