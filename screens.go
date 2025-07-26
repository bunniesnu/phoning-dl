package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (m *App) LoadingConfigScreen() *fyne.Container {
	title := widget.NewLabel("Loading configuration. Please wait...")
	defer func() {
		_ = m.LoadConfig()
	}()
	return container.NewVBox(
		title,
	)
}