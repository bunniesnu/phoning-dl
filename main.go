package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Initialize
	a := app.NewWithID(AppID)
	w := a.NewWindow(DefaultWindowTitle)
	size := fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight)
	w.Resize(size)
	w.CenterOnScreen()
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
			w.Resize(size)
			screen = m.MainScreen()
			w.SetContent(screen)
		})
	}()

	// Run
	w.ShowAndRun()
}