package main

import (
	"time"

	"fyne.io/fyne/v2"
)

type Config struct {
	ApiKey       string `json:"api_key"`
	AccessToken  string `json:"access_token"`
	TokenTimeout int64  `json:"token_timeout"`
}

type App struct {
	w      *fyne.Window
	config *Config
}

type LiveJSON struct {
	Id				  int    `json:"liveId"`
	Title			  string `json:"title"`
	MediaType		  string `json:"mediaType"`
	StartAt			  string `json:"startAt"`
	EndAt			  string `json:"endAt"`
	Duration		  int    `json:"duration"`
	ScreenOrientation string `json:"screenOrientation"`
}

type Live struct {
	Title 	    string
	Id		    int
	IsVideo     bool
	StartAt     time.Time
	EndAt       time.Time
	Duration    time.Duration
	IsLandscape bool
	Selected	bool
}