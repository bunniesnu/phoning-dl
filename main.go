package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Initialize
	a := app.New()
	w := a.NewWindow(DefaultWindowTitle)
	size := fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight)
	w.Resize(size)
	m := &App{
		w:      &w,
		config: &Config{},
	}

	// Contents
	err := m.UpdateState(0)
	if err != nil {
		w.SetContent(container.NewVBox(widget.NewLabel("Oops, something went wrong!")))
	}

	// Run
	w.ShowAndRun()
}