package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (m *App) DownloadScreen(liveSelection *[]Live, livesData *[]LiveJSON) fyne.CanvasObject {
	label := widget.NewLabel(fmt.Sprintf("Selected: %d / %d", getSelectedNum(liveSelection), len(*livesData)))
	refreshLabel := func() {
		label.SetText(fmt.Sprintf("Selected: %d / %d", getSelectedNum(liveSelection), len(*livesData)))
		label.Refresh()
	}
	list, checks := DrawList(liveSelection, InnerHeight, refreshLabel)
	header := container.NewHBox(
		label,
		layout.NewSpacer(),
		widget.NewButton("Select All", func() {
			for i := range *liveSelection {
				(*liveSelection)[i].Selected = true
				checks[i].SetChecked(true)
			}
			refreshLabel()
		}),
		widget.NewButton("Deselect All", func() {
			for i := range *liveSelection {
				(*liveSelection)[i].Selected = false
				checks[i].SetChecked(false)
			}
			refreshLabel()
		}),
	)
	vbox := container.NewVBox(
		header,
		list,
	)
	return vbox
}