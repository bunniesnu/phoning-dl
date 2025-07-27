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
	PNXMLInfo	*PNXMLInfo
}

type MetaData struct {
	Bitrate int
	FPS 	string
	Codec 	string
	Width 	int
	Height 	int
	URL 	string
	Size 	int64
}

type PNXMLInfo struct {
	MaxHeight int
	MetaDatas []MetaData
	ImageURL  string
}

type Representation struct {
	BaseURL []struct {
		Value string `json:"value"`
	}
	Width int `json:"width"`
	Height int `json:"height"`
	BandWidth int `json:"bandwidth"`
	FrameRate string `json:"frameRate"`
	Codec string `json:"codecs"`
}

type PNXMLJSON struct {
	Period []struct {
		AdaptationSet []struct {
			MaxHeight    float64 `json:"maxHeight"`
			MimeType     string  `json:"mimeType"`
			Representation []Representation `json:"representation"`
		} `json:"adaptationSet"`
		SupplementalProperty []struct {
			Any []struct {
				Cover []struct {
					Value string `json:"value"`
				} `json:"cover,omitempty"`
			} `json:"any"`
		} `json:"supplementalProperty"`
	} `json:"period"`
}