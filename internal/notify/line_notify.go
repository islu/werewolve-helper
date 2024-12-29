package notify

// IMPORTANT: Announcing the end of service for LINE Notify
// Ref: https://notify-bot.line.me/closing-announce

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const lineNotifyApiURL = "https://notify-api.line.me/api/notify"

// Deprecated: Send line notify simple text
func SendText(accessToken, message string) (err error) {
	body := strings.NewReader(url.Values{
		"message": []string{message},
	}.Encode())
	contentType := "application/x-www-form-urlencoded"

	err = sendToLineServer(body, accessToken, contentType)
	return
}

func sendToLineServer(body io.Reader, accessToken, contentType string) (err error) {
	req, err := http.NewRequest(
		"POST",
		lineNotifyApiURL,
		body,
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	var responseBody struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return
	}
	defer res.Body.Close()

	if responseBody.Status != 200 {
		err = fmt.Errorf("%d: %s", responseBody.Status, responseBody.Message)
	}
	return
}
