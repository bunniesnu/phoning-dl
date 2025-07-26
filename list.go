package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func DrawList(lives *[]Live, height float32, refresh func()) *container.Scroll {
	items := make([]fyne.CanvasObject, len(*lives))
	for i := 0; i < len(*lives); i++ {
		checkBox := widget.NewCheck(strconv.Itoa((*lives)[i].Id),
			func(b bool) {
				(*lives)[i].Selected = b
				if refresh != nil {
					refresh()
				}
			},
		)
		checkBox.SetChecked((*lives)[i].Selected)
		item := container.NewHBox(
			checkBox,
			widget.NewLabel((*lives)[i].Title),
			widget.NewLabel((*lives)[i].StartAt.Format("15:04:05")),
			widget.NewLabel((*lives)[i].EndAt.Format("15:04:05")),
			widget.NewLabel((*lives)[i].Duration.String()),
		)
		items[i] = item
	}
	vbox := container.NewVBox(items...)
	scrollable := container.NewVScroll(vbox)
	scrollable.SetMinSize(fyne.NewSize(0, height))
	return scrollable
}