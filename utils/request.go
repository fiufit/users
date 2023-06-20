package utils

import (
	"io"
	"net/http"
	"time"
)

func MakeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	return res, err
}
