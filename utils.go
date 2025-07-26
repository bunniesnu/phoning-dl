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
	"strings"
	"time"

	"github.com/bunniesnu/go-gmailnator"
	"github.com/bunniesnu/weverse-api"
	"github.com/chromedp/chromedp"
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

func GenerateAccessToken(updateProgress func(msg string, value float64)) (string, error) {
	// Generate a random email
	gmail, err := gmailnator.NewGmailnator()
	if err != nil {
		return "", fmt.Errorf("error creating Gmailnator client: %v", err)
	}
	err = gmail.GenerateEmail()
	if err != nil {
		return "", fmt.Errorf("error generating email: %v", err)
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
		return "", fmt.Errorf("error creating Weverse client: %v", err)
	}
	nickname, err := w.GetAccountNicknameSuggestion()
	if err != nil {
		return "", fmt.Errorf("error getting nickname suggestion")
	}
	w.Nickname = nickname
	err = w.CreateAccount()
	updateProgress("Signed up with Weverse", 0.3)
	if err != nil {
		return "", fmt.Errorf("error signing up: %v", err)
	}
	res := ""
	for i := range 5 {
		email, err := gmail.GetMails()
		if err != nil {
			return "", fmt.Errorf("error getting emails: %v", err)
		}
		if email == nil {
			return "", fmt.Errorf("no emails found")
		}
		for _, mail := range email {
			messageId := mail.Mid
			mailDetails, err := gmail.GetMailBody(messageId)
			if err != nil {
				return "", fmt.Errorf("error getting mail body for message ID %s: %v", messageId, err)
			}
			if mailDetails == "" {
				return "", fmt.Errorf("mail body is empty for message ID %s", messageId)
			}
			if strings.Contains(mailDetails, "account.weverse.io/signup") {
				start := strings.Index(mailDetails, "https://account.weverse.io/signup")
				if start == -1 {
					return "", fmt.Errorf("verification link not found in mail body")
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
		time.Sleep(5 * time.Second)
	}
	if res == "" {
		return "", fmt.Errorf("verification link not found in any emails")
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
		return "", fmt.Errorf("error clicking link: %v", err)
	}
	updateProgress("Clicked verification link", 0.6)

	// Check if the email is verified
	val, err := w.GetAccountStatus()
	if err != nil {
		return "", fmt.Errorf("error checking verification: %v", err)
	}
	if !(val.EmailVerified) {
		return "", fmt.Errorf("email verification failed")
	}
	updateProgress("Email verified successfully", 0.8)

	// Register the account to get the access token
	if email == "" || password == "" {
		return "", fmt.Errorf("Email or password not found in registration response")
	}
	body := map[string]any{
		"email":    email,
		"password": password,
	}
	encodedBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://sdk.weverse.io/api/v2/auth/token/by-credentials", strings.NewReader(string(encodedBody)))
	var DefaultHeaders = map[string]string{
		"Host": "sdk.weverse.io",
		"Accept": "*/*",
		"X-SDK-SERVICE-ID": "phoning",
		"X-SDK-LANGUAGE": "ko",
		"X-CLOG-USER-DEVICE-ID": "1",
		"X-SDK-PLATFORM": "iOS",
		"Accept-Language": "ko-KR,ko;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Content-Type": "application/json",
		"X-SDK-VERSION": "3.4.2",
		"User-Agent": "Phoning/20201014 CFNetwork/3826.500.131 Darwin/24.5.0",
		"Connection": "keep-alive",
		"X-SDK-TRACE-ID": "1",
		"X-SDK-APP-VERSION": "2.2.1",
		"Pragma": "no-cache",
		"Cache-Control": "no-cache",
		"X-SDK-SERVICE-SECRET": os.Getenv("PHONING_SDK_SERVICE_SECRET"),
	}
	for key, value := range DefaultHeaders {
		req.Header.Set(key, value)
	}
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
	updateProgress("Access token generated successfully", 1.0)
	return accessToken, nil
}