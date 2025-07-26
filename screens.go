package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoadingConfigScreen() *fyne.Container {
	title := widget.NewLabel("Loading configuration...")
	return container.NewVBox(
		title,
	)
}