package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TableColNum = 5

func DrawList(lives *[]Live, height float32, refresh func()) *container.Scroll {
	vbox := container.NewGridWithColumns(TableColNum)
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
		vbox.Add(checkBox)
		vbox.Add(widget.NewLabel((*lives)[i].Title))
		vbox.Add(widget.NewLabel((*lives)[i].StartAt.Format("15:04:05")))
		vbox.Add(widget.NewLabel((*lives)[i].EndAt.Format("15:04:05")))
		vbox.Add(widget.NewLabel((*lives)[i].Duration.String()))
	}
	scrollable := container.NewVScroll(vbox)
	scrollable.SetMinSize(fyne.NewSize(0, height))
	return scrollable
}