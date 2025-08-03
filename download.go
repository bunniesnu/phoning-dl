package main

import (
	"fmt"
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
	list, checks := DrawList(liveSelection, ListHeight, refreshLabel)
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
	var downloadFolder string
	downloadFolderLabel := widget.NewLabel("Download Folder: Not selected")
	downloadFolderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, *m.w)
			return
		}
		if uri != nil {
			slog.Info("Download folder selected", "path", uri.Path())
			downloadFolder = uri.Path()
			downloadFolderLabel.SetText(fmt.Sprintf("Download Folder: %s", downloadFolder))
			downloadFolderLabel.Refresh()
		}
	}, *m.w)
	downloadFolderPrompt := widget.NewButton("Select Download Folder", func() {
		slog.Info("Opening download folder dialog")
		downloadFolderDialog.Show()
	})
	footer := container.NewHBox(
		totalSizeLabel,
		layout.NewSpacer(),
		downloadFolderLabel,
		downloadFolderPrompt,
		widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
			if downloadFolder == "" {
				dialog.ShowError(fmt.Errorf("Please select a download folder"), *m.w)
				return
			}
			if info, err := os.Stat(downloadFolder); err != nil || !info.IsDir() {
				dialog.ShowError(fmt.Errorf("Invalid download folder selected"), *m.w)
				return
			}
			slog.Info("Starting download")
			go fyne.Do(func() { m.StartDownload(liveSelection, downloadFolder) })
		}),
	)
	vbox := container.NewVBox(
		header,
		list,
		footer,
	)
	return vbox
}

func (m *App) StartDownload(liveSelection *[]Live, baseDir string) {
	w := fyne.CurrentApp().NewWindow("PhoningDL - Download")
	w.Resize(fyne.NewSize(DownloadWindowWidth, DownloadWindowHeight))
	selNum, totalSize := getSelectedNum(liveSelection)
	totalProgress := widget.NewProgressBar()
	totalProgress.Max = float64(totalSize)
	completed := 0
	progressLabel := widget.NewLabel(fmt.Sprintf("Downloading (%d / %d)", completed, selNum))
	update := func() {
		completed++
		fyne.Do(func() { progressLabel.SetText(fmt.Sprintf("Downloading (%d / %d)", completed, selNum)) })
	}
	onProgress := func(progressDelta int64) {
		fyne.Do(func() { totalProgress.SetValue(totalProgress.Value + float64(progressDelta)) })
	}
	detailsVbox, cancel := DownloadConcurrentList(liveSelection, update, baseDir, onProgress)
	vbox := container.NewVBox(
		progressLabel,
		totalProgress,
		detailsVbox,
	)
	w.SetContent(vbox)
	w.SetCloseIntercept(func() {
		slog.Info("User requested download cancel")
		dialog.ShowConfirm("Confirm", "Are you sure you want to cancel the download?", func(confirm bool) {
			if confirm {
				slog.Info("User confirmed download cancel")
				cancel()
				w.Close()
			}
		}, w)
	})
	w.Show()
}
