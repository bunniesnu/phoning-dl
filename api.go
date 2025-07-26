package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func hash(url, apikey string) map[string]string {
	apiKey := []byte(apikey)

	msgpad := int(time.Now().UnixNano() / int64(time.Millisecond))

	if len(url) > 255 {
		url = url[:255]
	}

	message := []byte(url + strconv.Itoa(msgpad))

	mac := hmac.New(sha1.New, apiKey)
	mac.Write(message)
	digest := mac.Sum(nil)

	md := base64.StdEncoding.EncodeToString(digest)
	return map[string]string{
		"msgpad": strconv.Itoa(msgpad),
		"md": md,
	}
}

func getAPIHeaders(accessToken string) map[string]string {
	header := make(map[string]string)
	maps.Copy(header, DefaultAPIHeaders)
	if accessToken != "" {
		header["Authorization"] = "Bearer " + accessToken
	}
	return header
}

func Phoning(method, apiKey, accessToken, endpoint string, params ...map[string]string) (map[string]any, error) {
	var paramMap map[string]string
	if len(params) > 0 && params[0] != nil && method == "GET" {
		paramMap = params[0]
	} else {
		paramMap = make(map[string]string)
	}
	values := url.Values{}
	for k, v := range paramMap {
		values.Set(k, v)
	}
	encodeUrl := "https://apis.naver.com/phoning/phoning-api/api" + endpoint
	if len(values) > 0 {
		encodeUrl += "?" + values.Encode()
	}
	h := hash(encodeUrl, apiKey)
	hashValues := url.Values{}
	hashValues.Set("msgpad", h["msgpad"])
	hashValues.Set("md", h["md"])
	queryUrl := encodeUrl
	if len(values) > 0 {
		queryUrl += "&" + hashValues.Encode()
	} else {
		queryUrl += "?" + hashValues.Encode()
	}
	req, err := http.NewRequest(method, queryUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range getAPIHeaders(accessToken) {
		req.Header.Set(key, value)
	}
	body := make([]byte, 0)
	if method == "POST" || method == "PUT" {
		body, _ = json.Marshal(params[0])
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("phoning API request failed with status code %d", resp.StatusCode)
	}
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var response map[string]any
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode phoning API response: %w, %s", err, string(respBody))
	}
	return response, nil
}