package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/j-rooft/api/controllers"
	"github.com/j-rooft/api/models"
	. "github.com/j-rooft/api/tests"
)

func TestTelegramControllerSendMessage(t *testing.T) {
	testsData := []struct {
		description     string
		requestBody     string
		defaultChatID   *string
		outboundMessage models.OutboundTelegramMessage
		responseCode    int
	}{
		{
			description: "happy path",
			requestBody: `{"message":"opossum","chat_id":1234}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description:     "no message",
			requestBody:     `{"text":"opossum","chat_id":1234}`,
			outboundMessage: models.NewOutboundTelegramMessage(intPtr(1234)),
			responseCode:    http.StatusInternalServerError,
		},
		{
			description:     "empty message",
			requestBody:     `{"message":"","chat_id":1234}`,
			outboundMessage: models.NewOutboundTelegramMessage(intPtr(1234)),
			responseCode:    http.StatusInternalServerError,
		},
		{
			// everything is fine, but HTTP client will return error
			description: "telegram server error",
			requestBody: `{"message":"opossum","chat_id":1234}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			description: "no telegram chat ID",
			requestBody: `{"message":"opossum"}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              nil,
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			description:   "only default chat ID",
			requestBody:   `{"message":"opossum"}`,
			defaultChatID: strPtr("1111"),
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1111),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description: "default chat_id override",
			requestBody: `{"message":"opossum","chat_id":1111,"silent":true}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1111),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			defaultChatID: strPtr("2222"),
			responseCode:  http.StatusOK,
		},
		{
			description: "silent message",
			requestBody: `{"message":"opossum","chat_id":1234,"silent":true}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description: "non-silent message",
			requestBody: `{"message":"opossum","chat_id":1234,"silent":false}`,
			outboundMessage: models.OutboundTelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: false,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData)

		telegramController := TelegramController{
			Config: NewConfigSafe(strPtr("8080"), strPtr("1"), testData.defaultChatID, nil),
			HTTPClient: &MockHTTPClient{
				do: func(req *http.Request) (*http.Response, error) {
					m := models.NewOutboundTelegramMessage(nil)
					err := json.NewDecoder(req.Body).Decode(&m)
					if err != nil {
						panic(fmt.Sprintf("Cannot decode outbound telegram message: %s", err))
					}

					assert.Equal(testData.outboundMessage, m, "Outbound message is not as expected")

					var respErr error
					if testData.responseCode != http.StatusOK {
						respErr = errors.New("test error")
					}

					w := httptest.NewRecorder()
					w.WriteHeader(testData.responseCode)
					return w.Result(), respErr
				},
			},
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", bytes.NewReader([]byte(testData.requestBody)))

		telegramController.SendMessage(w, req)

		assert.Equal(testData.responseCode, w.Code, fmt.Sprintf("Response code is not correct: %s", formatBody(w)))
	}
}

func formatBody(w *httptest.ResponseRecorder) string {
	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		panic(fmt.Sprintf("Error while reading body: %s", err))
	}
	return string(bodyBytes)
}

type MockHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func strPtr(str string) *string {
	return &str
}

func intPtr(num int) *int {
	return &num
}
