package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const TableColNum = 5
const DateTimeFormat = "2006-01-02 15:04:05"

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
		vbox.Add(widget.NewLabel(live.StartAt.Format(DateTimeFormat)))
		vbox.Add(widget.NewLabel(live.EndAt.Format(DateTimeFormat)))
		durationHours := int(live.Duration.Hours())
		durationMinutes := int(live.Duration.Minutes()) % 60
		durationSeconds := int(live.Duration.Seconds()) % 60
		if durationHours == 0 {
			vbox.Add(widget.NewLabel(fmt.Sprintf("%02d:%02d", durationMinutes, durationSeconds)))
		} else {
			vbox.Add(widget.NewLabel(fmt.Sprintf("%d:%02d:%02d", durationHours, durationMinutes, durationSeconds)))
		}
	}
	scrollable := container.NewVScroll(vbox)
	scrollable.SetMinSize(fyne.NewSize(0, height))
	return scrollable, checks
}