package main

import (
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (m *App) LoadingConfigScreen(done chan struct{}) *fyne.Container {
	title := widget.NewLabel("Loading configuration. Please wait...")
	progress := widget.NewProgressBar()
	updateProgress := func(msg string, value float64) {
		slog.Info(msg)
		fyne.Do(func() {
			progress.SetValue(value)
		})
	}
	vbox := container.NewVBox(
		title,
		progress,
	)
	var retryBtn *widget.Button
	validateConfig := func() {
		fyne.Do(func() {
			progress.Show()
			progress.SetValue(0.0)
		})
		err := m.LoadConfig()
		if err != nil || m.config.AccessToken == "" || m.config.TokenTimeout < time.Now().Unix() {
			accessToken, expiresIn, err := GenerateAccessToken(updateProgress)
			if err != nil {
				slog.Error("Failed to generate access token", "error", err)
				vbox.Add(widget.NewLabel("Failed to generate access token."))
				vbox.Add(retryBtn)
				return
			}
			m.config.AccessToken = accessToken
			m.config.TokenTimeout = time.Now().Unix() + expiresIn
			err = m.SaveConfig()
			if err != nil {
				slog.Error("Failed to save access token", "error", err)
				vbox.Add(widget.NewLabel("Failed to save configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		if m.config.ApiKey == "" {
			m.config.ApiKey = os.Getenv("PHONING_API_KEY")
			err := m.SaveConfig()
			if err != nil {
				slog.Error("Failed to save API key", "error", err)
				vbox.Add(widget.NewLabel("Failed to save configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		updateProgress("Validating configuration", 0.9)
		_, err = Phoning("GET", m.config.ApiKey, m.config.AccessToken, "/fan/v1.0/users/me")
		if err != nil {
			slog.Info("Trying to login")
			_, err = Phoning("POST", m.config.ApiKey, "", "/fan/v1.0/login", map[string]string{
				"wevAccessToken": m.config.AccessToken,
				"tokenType": "APNS",
				"deviceToken": "",
			})
			if err != nil {
				slog.Error("Failed to login", "error", err)
				vbox.Add(widget.NewLabel("Failed to login."))
				vbox.Add(retryBtn)
				return
			}
			slog.Info("Login successful")
			_, err = Phoning("GET", m.config.ApiKey, m.config.AccessToken, "/fan/v1.0/users/me")
			if err != nil {
				slog.Error("Failed to validate configuration", "error", err)
				vbox.Add(widget.NewLabel("Failed to validate configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		updateProgress("Configuration validated successfully", 1.0)
		fyne.Do(func() {
			progress.Hide()
		})
		vbox.Add(widget.NewLabel("Configuration loaded successfully!"))
		done <- struct{}{}
	}
	retryBtn = widget.NewButton("Retry", func() {
		vbox.RemoveAll()
		vbox.Add(title)
		vbox.Add(progress)
		go validateConfig()
	})
	go validateConfig()
	return vbox
}

func (m *App) MainScreen() *fyne.Container {
	title := widget.NewLabel("PhoningDL")
	vbox := container.NewVBox(
		title,
	)
	return vbox
}