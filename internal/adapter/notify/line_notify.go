// IMPORTANT: Announcing the end of service for LINE Notify.
// Refer to https://notify-bot.line.me/closing-announce.

package notify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const lineNotifyURL = "https://notify-api.line.me/api/notify"

// Deprecated: end of service for LINE Notify
//
// Send line notify simple text
func SendText(accessToken, message string) error {
	// Set the request body with the message
	body := strings.NewReader(url.Values{
		"message": []string{message},
	}.Encode())
	contentType := "application/x-www-form-urlencoded"

	// Send the request to the LINE Notify API
	err := requestToLineServer(body, accessToken, contentType)
	return err
}

// requestToLineServer handles the HTTP request to the LINE Notify API.
// It takes the request body, access token, and content type as input.
// It returns an error if the request fails or if the API returns a non-200 status.
func requestToLineServer(body io.Reader, accessToken, contentType string) error {
	// Create a new HTTP POST request with the provided body and URL
	req, err := http.NewRequest("POST", lineNotifyURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// Send the request to the LINE Notify API
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	// Define a struct to hold the response body
	var responseBody struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if responseBody.Status != 200 {
		err := fmt.Errorf("%d: %s", responseBody.Status, responseBody.Message)
		return err
	}
	return nil
}
