package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bunniesnu/go-gmailnator"
	"github.com/bunniesnu/weverse-api"
	"github.com/chromedp/chromedp"
	"golang.org/x/sync/errgroup"
)

const (
	lower       = "abcdefghijklmnopqrstuvwxyz"
	upper       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits      = "0123456789"
	specials    = "!@#%^_=+"
	allChars    = lower + upper + digits + specials
	passwordLen = 16
)

func getRandomChar(charset string) byte {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	return charset[n.Int64()]
}

func GenerateAccessToken(updateProgress func(msg string, value float64)) (string, int64, error) {
	// Generate a random email
	gmail, err := gmailnator.NewGmailnator()
	if err != nil {
		return "", 0, fmt.Errorf("error creating Gmailnator client: %v", err)
	}
	err = gmail.GenerateEmail()
	if err != nil {
		return "", 0, fmt.Errorf("error generating email: %v", err)
	}
	email := gmail.Email.Email
	updateProgress("Generated random email", 0.1)

	// Generate a random password
	passwordSet := []byte{
		getRandomChar(lower),
		getRandomChar(upper),
		getRandomChar(digits),
		getRandomChar(specials),
	}
	for i := 4; i < 16; i++ {
		passwordSet = append(passwordSet, getRandomChar(allChars))
	}
	for i := len(passwordSet) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		passwordSet[i], passwordSet[j.Int64()] = passwordSet[j.Int64()], passwordSet[i]
	}
	password := string(passwordSet)
	updateProgress("Generated random password", 0.2)

	// Create a Weverse client and sign up
	w, err := weverse.New(email, password, "", 0)
	if err != nil {
		return "", 0, fmt.Errorf("error creating Weverse client: %v", err)
	}
	nickname, err := w.GetAccountNicknameSuggestion()
	if err != nil {
		return "", 0, fmt.Errorf("error getting nickname suggestion")
	}
	w.Nickname = nickname
	err = w.CreateAccount()
	updateProgress("Signed up with Weverse", 0.3)
	if err != nil {
		return "", 0, fmt.Errorf("error signing up: %v", err)
	}
	res := ""
	for i := range 5 {
		time.Sleep(5 * time.Second)
		email, err := gmail.GetMails()
		if err != nil {
			return "", 0, fmt.Errorf("error getting emails: %v", err)
		}
		if email == nil {
			return "", 0, fmt.Errorf("no emails found")
		}
		for _, mail := range email {
			messageId := mail.Mid
			mailDetails, err := gmail.GetMailBody(messageId)
			if err != nil {
				return "", 0, fmt.Errorf("error getting mail body for message ID %s: %v", messageId, err)
			}
			if mailDetails == "" {
				return "", 0, fmt.Errorf("mail body is empty for message ID %s", messageId)
			}
			if strings.Contains(mailDetails, "account.weverse.io/signup") {
				start := strings.Index(mailDetails, "https://account.weverse.io/signup")
				if start == -1 {
					return "", 0, fmt.Errorf("verification link not found in mail body")
				}
				end := strings.IndexAny(mailDetails[start:], " \"'<")
				if end == -1 {
					res = mailDetails[start:]
				} else {
					res = mailDetails[start : start+end]
				}
				break
			}
		}
		if res != "" {
			break
		}
		updateProgress(fmt.Sprintf("Checking for verification email (%d/5)", i+1), 0.3+float64(i+1)*0.02)
	}
	if res == "" {
		return "", 0, fmt.Errorf("verification link not found in any emails")
	}
	updateProgress("Found verification link", 0.5)
	res = strings.ReplaceAll(res, "&amp;", "&")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()
	var html string
	err = chromedp.Run(ctx,
		chromedp.Navigate(res),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return "", 0, fmt.Errorf("error clicking link: %v", err)
	}
	updateProgress("Clicked verification link", 0.6)

	// Check if the email is verified
	time.Sleep(3 * time.Second)
	val, err := w.GetAccountStatus()
	if err != nil {
		return "", 0, fmt.Errorf("error checking verification: %v", err)
	}
	if !(val.EmailVerified) {
		return "", 0, fmt.Errorf("email verification failed")
	}
	updateProgress("Email verified successfully", 0.7)

	// Register the account to get the access token
	if email == "" || password == "" {
		return "", 0, fmt.Errorf("Email or password not found in registration response")
	}
	body := map[string]any{
		"email":    email,
		"password": password,
	}
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return "", 0, err
	}
	req, err := http.NewRequest("POST", "https://sdk.weverse.io/api/v2/auth/token/by-credentials", strings.NewReader(string(encodedBody)))
	for key, value := range DefaultWevSDKHeaders {
		req.Header.Set(key, value)
	}
	req.Header.Set("X-SDK-SERVICE-SECRET", os.Getenv("PHONING_SDK_SERVICE_SECRET"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	decodedResponse := make(map[string]any)
	respBody, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &decodedResponse); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}
	accessToken, ok := decodedResponse["accessToken"].(string)
	if !ok {
		log.Fatal("Access token not found in response")
	}
	expiresIn, ok := decodedResponse["expiresIn"].(float64)
	if !ok {
		log.Fatal("Expires in not found in response")
	}
	updateProgress("Access token generated successfully", 0.8)
	return accessToken, int64(expiresIn), nil
}

func (m *App) FetchLives() (*[]LiveJSON, error) {
	var callsData *[]LiveJSON = new([]LiveJSON)
	nextCursor := ""
	cnt := 0
	for {
		cnt++
		if cnt > 10 {
			return nil, fmt.Errorf("too many iterations, stopping to prevent infinite loop")
		}
		params := map[string]string{"limit": "100"}
		if nextCursor != "" {
			params["cursor"] = nextCursor
		}
		calls, err := Phoning("GET", m.config.ApiKey, m.config.AccessToken, "/fan/v1.0/lives", params)
		if err != nil {
			log.Fatalf("%v", err)
		}
		rawData, ok := calls["data"].([]any)
		if !ok {
			log.Fatalf("Unexpected data format: %T", calls["data"])
		}
		for _, item := range rawData {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				log.Fatalf("Error marshaling item: %v", err)
			}
			var live LiveJSON
			if err := json.Unmarshal(itemBytes, &live); err != nil {
				log.Fatalf("Error unmarshaling item to LiveJSON: %v", err)
			}
			*callsData = append(*callsData, live)
		}
		cursors, ok := calls["cursors"].(map[string]any)
		if !ok {
			log.Fatalf("Unexpected cursors format: %T", calls["cursors"])
		}
		next, ok := cursors["next"].(string)
		if !ok {
			break
		}
		nextCursor = next
	}
	return callsData, nil
}

func getSelectedNum(liveSelection *[]Live) (int, int64) {
	selectedCount := 0
	totalSize := int64(0)
	for _, live := range *liveSelection {
		if live.Selected {
			selectedCount++
			for _, metaData := range live.PNXMLInfo.MetaDatas {
				if metaData.Height == live.SelHeight {
					totalSize += metaData.Size
					break
				}
			}
		}
	}
	return selectedCount, totalSize
}

func getFileSizeWithContext(url string, ctx context.Context) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making HEAD request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	return resp.ContentLength, nil
}

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

// safeCreateFile ensures destPath is within baseDir and that there is enough free space.
func safeCreateFile(destPath, baseDir string, size int64) (*os.File, error) {
	cleanDest := filepath.Clean(destPath)

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("resolving base dir: %w", err)
	}
	absDest, err := filepath.Abs(cleanDest)
	if err != nil {
		return nil, fmt.Errorf("resolving dest path: %w", err)
	}

	rel, err := filepath.Rel(absBase, absDest)
	if err != nil {
		return nil, fmt.Errorf("resolving relative path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return nil, fmt.Errorf("destination %q is outside of %q", absDest, absBase)
	}

	freeBytes, err := getDiskFreeSpace(absBase)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if freeBytes < size {
		return nil, fmt.Errorf("not enough disk space: need %d, have %d", size, freeBytes)
	}

	outFile, err := os.Create(absDest)
	if err != nil {
		return nil, fmt.Errorf("creating file: %w", err)
	}
	if err := outFile.Truncate(size); err != nil {
		outFile.Close()
		return nil, fmt.Errorf("preallocating file: %w", err)
	}
	return outFile, nil
}

func DownloadVideo(
	ctx context.Context,
	url, destPath, baseDir string,
	concurrency int,
	updateFunc func(value int64),
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return fmt.Errorf("creating HEAD request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HEAD request failed: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HEAD returned %s", resp.Status)
	}
	length := resp.ContentLength
	if length <= 0 || length > MaxAllowedSize {
		return fmt.Errorf("invalid or too large content length: %d", length)
	}
	supportRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	outFile, err := safeCreateFile(destPath, baseDir, length)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var downloaded int64
	var progressMu sync.Mutex

	updateProgress := func(n int64) {
		progressMu.Lock()
		updateFunc(atomic.AddInt64(&downloaded, n))
		progressMu.Unlock()
	}

	if !supportRanges || concurrency <= 1 {
		return singleDownload(ctx, url, outFile, updateProgress)
	}

	partSize := length / int64(concurrency)
	eg, ctx := errgroup.WithContext(ctx)

	for i := range concurrency {
		start := int64(i) * partSize
		end := start + partSize - 1
		if i == concurrency-1 {
			end = length - 1
		}
		chunkStart, chunkEnd := start, end

		eg.Go(func() error {
			var lastErr error
			for attempt := range MaxRetries {
				if err := downloadChunk(ctx, url, outFile, chunkStart, chunkEnd, updateProgress); err != nil {
					lastErr = err
					time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
					continue
				}
				return nil
			}
			return fmt.Errorf("chunk %d-%d failed: %w", chunkStart, chunkEnd, lastErr)
		})
	}
	return eg.Wait()
}

func singleDownload(ctx context.Context, url string, outFile *os.File, updateProgress func(int64)) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := outFile.Write(buf[:n]); err != nil {
				return err
			}
			updateProgress(int64(n))
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	return nil
}

func downloadChunk(ctx context.Context, url string, outFile *os.File, start, end int64, updateProgress func(int64)) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("expected 206 for range %d-%d, got %d", start, end, resp.StatusCode)
	}

	buf := make([]byte, 32*1024)
	offset := start
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := outFile.WriteAt(buf[:n], offset); err != nil {
				return err
			}
			offset += int64(n)
			updateProgress(int64(n))
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	return nil
}
