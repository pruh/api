package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/glog"
	"github.com/pruh/api/config"

	apihttp "github.com/pruh/api/http"
)

// Controller stores config and HTTP client for requests.
type Controller struct {
	Config     *config.Configuration
	HTTPClient apihttp.Client
}

// SendMessage sends a message to Telegram and returns Telegram's response.
func (c *Controller) SendMessage(w http.ResponseWriter, r *http.Request) {
	m := NewMessage(c.Config.DefaultChatID)
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		glog.Errorf("Cannot decode body. %s", err)
		http.Error(w, fmt.Sprintf("Cannot decode body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if m.ChatID == nil {
		glog.Errorf("ChatID not set. %s", err)
		http.Error(w, "ChatID not set", http.StatusInternalServerError)
		return
	}
	if m.Message == "" {
		glog.Errorf("Message should not be empty. %s", err)
		http.Error(w, "Message should not be empty", http.StatusInternalServerError)
		return
	}

	resp, err := sendTelegram(m.Message, m.ChatID, m.Silent, c.Config.TelegramBoToken, c.HTTPClient)
	if err != nil {
		glog.Errorf("Cannot send message to telegram. %s", err)
		http.Error(w, fmt.Sprintf("Cannot send message to telegram: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	glog.Infof("telegram response code: %d headers: %+v\n", resp.StatusCode, resp.Header)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

/**
 * Utility function to send message to Telegram using REST API.
 */
func sendTelegram(text string, chatID *int, silent bool, botToken *string, httpClient apihttp.Client) (*http.Response, error) {
	m := NewTelegramMessage(chatID)
	m.DisableNotification = silent
	m.Text = text

	jsonStr, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", *botToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	glog.Infof("sending message to telegram: %s", jsonStr)
	return httpClient.Do(req)
}
