package main

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (m *App) DownloadScreen(liveSelection *[]Live, livesData *[]LiveJSON) fyne.CanvasObject {
	selNum, totalSize := getSelectedNum(liveSelection)
	totalSizeLabel := widget.NewLabel(fmt.Sprintf("Total Size: %s", formatBytes(totalSize)))
	label := widget.NewLabel(fmt.Sprintf("Selected: %d / %d", selNum, len(*livesData)))
	refreshLabel := func() {
		selNum, totalSize := getSelectedNum(liveSelection)
		label.SetText(fmt.Sprintf("Selected: %d / %d", selNum, len(*livesData)))
		label.Refresh()
		totalSizeLabel.SetText(fmt.Sprintf("Total Size: %s", formatBytes(totalSize)))
		totalSizeLabel.Refresh()
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
	footer := container.NewHBox(
		totalSizeLabel,
		layout.NewSpacer(),
		widget.NewButton("Download", func() {
			slog.Info("Starting download")
		}),
	)
	vbox := container.NewVBox(
		header,
		list,
		footer,
	)
	return vbox
}