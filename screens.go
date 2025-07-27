package main

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (m *App) LoadingConfigScreen(done chan struct{}) *fyne.Container {
	title := widget.NewLabel("Loading configuration. Please wait...")
	progress := widget.NewProgressBar()
	updateProgress := func(msg string, value float64) {
		slog.Info(msg)
		fyne.Do(func() {
			progress.SetValue(value)
		})
	}
	vbox := container.NewVBox(
		title,
		progress,
	)
	var retryBtn *widget.Button
	validateConfig := func() {
		fyne.Do(func() {
			progress.Show()
			progress.SetValue(0.0)
		})
		err := m.LoadConfig()
		if err != nil || m.config.AccessToken == "" || m.config.TokenTimeout < time.Now().Unix() {
			accessToken, expiresIn, err := GenerateAccessToken(updateProgress)
			if err != nil {
				slog.Error("Failed to generate access token", "error", err)
				vbox.Add(widget.NewLabel("Failed to generate access token."))
				vbox.Add(retryBtn)
				return
			}
			m.config.AccessToken = accessToken
			m.config.TokenTimeout = time.Now().Unix() + expiresIn
			err = m.SaveConfig()
			if err != nil {
				slog.Error("Failed to save access token", "error", err)
				vbox.Add(widget.NewLabel("Failed to save configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		if m.config.ApiKey == "" {
			m.config.ApiKey = os.Getenv("PHONING_API_KEY")
			err := m.SaveConfig()
			if err != nil {
				slog.Error("Failed to save API key", "error", err)
				vbox.Add(widget.NewLabel("Failed to save configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		updateProgress("Validating configuration", 0.9)
		_, err = Phoning("GET", m.config.ApiKey, m.config.AccessToken, "/fan/v1.0/users/me")
		if err != nil {
			slog.Info("Trying to login")
			_, err = Phoning("POST", m.config.ApiKey, "", "/fan/v1.0/login", map[string]string{
				"wevAccessToken": m.config.AccessToken,
				"tokenType": "APNS",
				"deviceToken": "",
			})
			if err != nil {
				slog.Error("Failed to login", "error", err)
				vbox.Add(widget.NewLabel("Failed to login."))
				vbox.Add(retryBtn)
				return
			}
			slog.Info("Login successful")
			_, err = Phoning("GET", m.config.ApiKey, m.config.AccessToken, "/fan/v1.0/users/me")
			if err != nil {
				slog.Error("Failed to validate configuration", "error", err)
				vbox.Add(widget.NewLabel("Failed to validate configuration."))
				vbox.Add(retryBtn)
				return
			}
		}
		updateProgress("Configuration validated successfully", 1.0)
		fyne.Do(func() {
			progress.Hide()
		})
		vbox.Add(widget.NewLabel("Configuration loaded successfully!"))
		done <- struct{}{}
	}
	retryBtn = widget.NewButton("Retry", func() {
		vbox.RemoveAll()
		vbox.Add(title)
		vbox.Add(progress)
		go validateConfig()
	})
	go validateConfig()
	return vbox
}

func (m *App) MainScreen() *fyne.Container {
	liveSelection := new([]Live)
	slog.Info("Loading main screen")
	fetchingLiveLabel := widget.NewLabel("Fetching lives...")
	vbox := container.NewVBox(
		fetchingLiveLabel,
	)
	retryBtn := new(widget.Button)
	loadFunc := func() {
		lives, err := m.FetchLives()
		if err != nil {
			slog.Error("Failed to fetch lives", "error", err)
			fyne.Do(func() {
				vbox.RemoveAll()
				vbox.Add(widget.NewLabel("Failed to fetch lives."))
				vbox.Add(retryBtn)
				vbox.Refresh()
			})
			return
		}
		slog.Info("Lives fetched successfully", "count", len(*lives))
		cnt := 0
		progress := widget.NewProgressBar()
		vbox.Add(progress)
		decodeFailed := false
		parseFunc := func(live LiveJSON, ctx context.Context) (*Live, error) {
			startAtParse, err := time.Parse(time.RFC3339Nano, live.StartAt)
			if err != nil {
				slog.Error("Failed to parse startAt", "error", err, "startAt", live.StartAt)
				return nil, err
			}
			endAtParse, err := time.Parse(time.RFC3339Nano, live.EndAt)
			if err != nil {
				slog.Error("Failed to parse endAt", "error", err, "endAt", live.EndAt)
				return nil, err
			}
			pnxmlData, err := getPNXML(m.config.ApiKey, m.config.AccessToken, live.Id)
			if err != nil {
				slog.Error("Failed to get PNXML data", "error", err, "liveId", live.Id)
				return nil, err
			}
			liveReturn := new(Live)
			maxSize := 0
			for _, metaData := range pnxmlData.MetaDatas {
				if metaData.Height > maxSize {
					maxSize = metaData.Height
				}
			}
			*liveReturn = Live{
				Id:       live.Id,
				Title:    live.Title,
				Selected: true,
				SelHeight: maxSize,
				IsVideo:  live.MediaType == "LIVE",
				StartAt: startAtParse,
				EndAt:   endAtParse,
				Duration: time.Duration(live.Duration) * time.Millisecond,
				IsLandscape: live.ScreenOrientation == "LANDSCAPE",
				PNXMLInfo: pnxmlData,
			}
			fyne.Do(func(){
				cnt++
				progress.SetValue(float64(cnt) / float64(len(*lives)))
			})
			return liveReturn, nil
		}
		liveSelectionTmp, err := concurrentExecute(parseFunc, *lives, DefaultFetchConcurrency)
		if err != nil {
			slog.Error("Failed to decode lives", "error", err)
			decodeFailed = true
		}
		slog.Info("Lives decoded successfully", "count", len(liveSelectionTmp))
		*liveSelection = make([]Live, 0)
		for _, live := range liveSelectionTmp {
			if live != nil {
				*liveSelection = append(*liveSelection, *live)
			} else {
				slog.Warn("Received nil live from concurrent execution")
			}
		}
		sort.Slice(*liveSelection, func(i, j int) bool {
			return (*liveSelection)[i].StartAt.After((*liveSelection)[j].StartAt)
		})
		fyne.Do(func() {
			vbox.RemoveAll()
			if decodeFailed {
				vbox.Add(widget.NewLabel("Failed to decode lives."))
				vbox.Add(retryBtn)
			} else {
				downloadScreen := m.DownloadScreen(liveSelection, lives)
				vbox.Add(downloadScreen)
			}
			vbox.Refresh()
		})
	}
	retryBtn = widget.NewButton("Retry", func() {
		fyne.Do(func() {
			vbox.RemoveAll()
			vbox.Add(fetchingLiveLabel)
			vbox.Refresh()
		})
		loadFunc()
	})
	go loadFunc()
	return vbox
}