package main

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (m *App) LoadingConfigScreen() *fyne.Container {
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
	go func() {
		err := m.LoadConfig()
		if err != nil || m.config.AccessToken == "" {
			accessToken, err := GenerateAccessToken(updateProgress)
			if err != nil {
				slog.Error("Failed to generate access token", "error", err)
				vbox.Add(widget.NewLabel("Failed to generate access token. Please try again later."))
				return
			}
			m.config.AccessToken = accessToken
			err = m.SaveConfig()
			if err != nil {
				vbox.Add(widget.NewLabel("Failed to save configuration. Please try again later."))
				return
			}
		}
		progress.Hide()
		vbox.Add(widget.NewLabel("Configuration loaded successfully!"))
	}()
	return vbox
}