package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TableColNum = 5

func DrawList(lives *[]Live, height float32, refresh func()) (*container.Scroll, []*widget.Check) {
	vbox := container.NewGridWithColumns(TableColNum)
	checks := make([]*widget.Check, 0, len(*lives))
	for i := 0; i < len(*lives); i++ {
		live := &(*lives)[i]
		checkBox := widget.NewCheck(strconv.Itoa(live.Id),
			func(b bool) {
				live.Selected = b
				if refresh != nil {
					refresh()
				}
			},
		)
		checkBox.SetChecked(live.Selected)
		checks = append(checks, checkBox)
		vbox.Add(checkBox)
		vbox.Add(widget.NewLabel(live.Title))
		vbox.Add(widget.NewLabel(live.StartAt.Format("15:04:05")))
		vbox.Add(widget.NewLabel(live.EndAt.Format("15:04:05")))
		vbox.Add(widget.NewLabel(live.Duration.String()))
	}
	scrollable := container.NewVScroll(vbox)
	scrollable.SetMinSize(fyne.NewSize(0, height))
	return scrollable, checks
}