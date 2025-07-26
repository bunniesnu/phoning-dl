package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize
	a := app.New()
	w := a.NewWindow(DefaultWindowTitle)
	loadingSize := fyne.NewSize(LoadingWindowWidth, LoadingWindowHeight)
	mainSize := fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight)
	w.Resize(loadingSize)
	m := &App{
		w:      &w,
		config: &Config{},
	}

	// Loading screen
	done := make(chan struct{})
	screen := m.LoadingConfigScreen(done)
	w.SetContent(screen)

	// Main screen
	go func() {
		<-done
		time.Sleep(100 * time.Millisecond)
		fyne.Do(func() {
			w.SetContent(m.MainScreen())
			w.Resize(mainSize)
			screen = m.MainScreen()
			w.SetContent(screen)
		})
	}()

	// Run
	w.ShowAndRun()
}