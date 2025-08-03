package main

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func DownloadConcurrentList(liveList *[]Live, updateFunc func(), baseDir string, onProgress func(int64)) (fyne.CanvasObject, context.CancelFunc) {
	liveListFiltered := make([]*Live, 0, len(*liveList))
	for _, live := range *liveList {
		if live.Selected {
			liveListFiltered = append(liveListFiltered, &live)
		}
	}
	vbox := container.NewVBox()
	add := func(live *Live, ctx context.Context) (*fyne.Container, error) {
		// Initialize
		var selMetaData *MetaData
		for _, metaData := range live.PNXMLInfo.MetaDatas {
			if metaData.Height == live.SelHeight {
				selMetaData = &metaData
				break
			}
		}
		size := selMetaData.Size
		if size == 0 {
			slog.Error("No size found for live", "liveId", live.Id)
			return nil, fmt.Errorf("no size found for live %d", live.Id)
		}
		progress := widget.NewProgressBar()
		wrappedProgress := container.NewWithoutLayout(progress)
		progressLabel := widget.NewLabel(fmt.Sprintf("%d", live.Id))
		hbox := container.NewHBox(
			progressLabel,
			wrappedProgress,
		)
		progress.Resize(fyne.NewSize(DownloadWindowWidth*2/3, min(hbox.MinSize().Height, progress.MinSize().Height)))
		fyne.Do(func() {
			vbox.Add(hbox)
		})
		time.Sleep(100 * time.Millisecond)

		// Download
		destPath := path.Join(baseDir, fmt.Sprintf("%d.mp4", live.Id))
		prog := int64(0)
		updateProgress := func(progressValue int64) {
			fyne.Do(func() {
				progress.SetValue(float64(progressValue) / float64(size))
			})
			onProgress(progressValue - prog)
			prog = progressValue
		}
		err := DownloadVideo(ctx, selMetaData.URL, destPath, baseDir, 10, updateProgress)
		if err != nil {
			if ctx.Err() == context.Canceled {
				slog.Info(fmt.Sprintf("Download for id %d canceled", live.Id))
				return nil, err
			}
			slog.Error(fmt.Sprintf("%s", err))
			return nil, err
		}
		updateFunc()
		slog.Info(fmt.Sprintf("Download for id %d complete", live.Id))
		return hbox, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	go concurrentExecuteAnyWithContext(add, liveListFiltered, DefaultConcurrency, ctx)
	scrollable := container.NewScroll(vbox)
	scrollable.SetMinSize(fyne.NewSize(0, DownloadWindowHeight*3/4))
	return scrollable, cancel
}
