package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/j-rooft/api/models"
	"github.com/j-rooft/api/utils"
)

type TelegramController struct {
	Config     *utils.Configuration
	HTTPClient utils.HTTPClient
}

func (c *TelegramController) SendMessage(w http.ResponseWriter, r *http.Request) {
	m := models.NewInboundTelegramMessage(c.Config.DefaultChatID)
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("sending message to telegram: %+v\n", m)

	if m.Message == "" {
		http.Error(w, "Message should not be empty", http.StatusInternalServerError)
		return
	}

	resp, err := sendTelegram(m.Message, m.ChatID, m.Silent, c.Config.TelegramBoToken, c.HTTPClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot send message to telegram: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("telegram response code: %d headers: %+v\n", resp.StatusCode, resp.Header)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

/**
 * Utility function to send message to Telegram using REST API.
 */
func sendTelegram(text string, chatID *int, silent bool, botToken *string, httpClient utils.HTTPClient) (*http.Response, error) {
	m := models.NewOutboundTelegramMessage(chatID)
	m.DisableNotification = silent
	m.Text = text

	jsonStr, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot/%s/sendMessage", *botToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return httpClient.Do(req)
}
