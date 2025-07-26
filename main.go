package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize
	a := app.New()
	w := a.NewWindow(DefaultWindowTitle)
	size := fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight)
	w.Resize(size)

	// Contents
	w.SetContent(LoadingConfigScreen())

	// Run
	w.ShowAndRun()
}