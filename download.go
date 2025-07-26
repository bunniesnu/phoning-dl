package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (m *App) DownloadScreen(liveSelection *[]Live, livesData *[]LiveJSON) fyne.CanvasObject {
	label := widget.NewLabel(fmt.Sprintf("Selected: %d / %d", getSelectedNum(liveSelection), len(*livesData)))
	refreshLabel := func() {
		label.SetText(fmt.Sprintf("Selected: %d / %d", getSelectedNum(liveSelection), len(*livesData)))
		label.Refresh()
	}
	list := DrawList(liveSelection, InnerHeight, refreshLabel)
	vbox := container.NewVBox(
		label,
		list,
	)
	return vbox
}