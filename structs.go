package main

import "fyne.io/fyne/v2"

type Config struct {
	ApiKey       string `json:"api_key"`
	AccessToken  string `json:"access_token"`
	TokenTimeout int64  `json:"token_timeout"`
}

type App struct {
	w      *fyne.Window
	config *Config
}