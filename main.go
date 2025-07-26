package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize
	a := app.New()

	// Load configuration
	loadingWindow := a.NewWindow(DefaultWindowTitle)
	size := fyne.NewSize(LoadingWindowWidth, LoadingWindowHeight)
	loadingWindow.Resize(size)
	m := &App{
		w:      &loadingWindow,
		config: &Config{},
	}

	// Contents
	screen := m.LoadingConfigScreen()
	loadingWindow.SetContent(screen)

	// Run
	loadingWindow.ShowAndRun()
}