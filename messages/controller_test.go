package messages_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/pruh/api/config/tests"
	"github.com/pruh/api/messages"
	. "github.com/pruh/api/messages"
)

func TestTelegramControllerSendMessage(t *testing.T) {
	testsData := []struct {
		description             string
		requestBody             string
		defaultChatID           *string
		botToken                *string
		telegramShouldBeCalled  bool
		expectedOutboundMessage *messages.TelegramMessage
		responseCode            int
	}{
		{
			description:            "happy path",
			requestBody:            `{"message":"opossum","chat_id":1234}`,
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description:             "malformed json",
			requestBody:             `{"message":"opossum","chat_id":1234`,
			telegramShouldBeCalled:  false,
			expectedOutboundMessage: nil,
			responseCode:            http.StatusBadRequest,
		},
		{
			description:             "no message",
			requestBody:             `{"text":"opossum","chat_id":1234}`,
			telegramShouldBeCalled:  false,
			expectedOutboundMessage: nil,
			responseCode:            http.StatusBadRequest,
		},
		{
			description:             "empty message",
			requestBody:             `{"message":"","chat_id":1234}`,
			telegramShouldBeCalled:  false,
			expectedOutboundMessage: nil,
			responseCode:            http.StatusBadRequest,
		},
		{
			// everything is fine, but HTTP client will return error
			description:            "telegram server error",
			requestBody:            `{"message":"opossum","chat_id":1234}`,
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			description:             "no telegram chat ID",
			requestBody:             `{"message":"opossum"}`,
			telegramShouldBeCalled:  false,
			expectedOutboundMessage: nil,
			responseCode:            http.StatusBadRequest,
		},
		{
			description:            "only default chat ID",
			requestBody:            `{"message":"opossum"}`,
			defaultChatID:          strPtr("1111"),
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1111),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description:            "default chat_id override",
			requestBody:            `{"message":"opossum","chat_id":1111,"silent":true}`,
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1111),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			defaultChatID: strPtr("2222"),
			responseCode:  http.StatusOK,
		},
		{
			description:            "silent message",
			requestBody:            `{"message":"opossum","chat_id":1234,"silent":true}`,
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: true,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description:            "non-silent message",
			requestBody:            `{"message":"opossum","chat_id":1234,"silent":false}`,
			telegramShouldBeCalled: true,
			expectedOutboundMessage: &messages.TelegramMessage{
				ChatID:              intPtr(1234),
				DisablePreview:      true,
				DisableNotification: false,
				Text:                "opossum",
			},
			responseCode: http.StatusOK,
		},
		{
			description:             "invalid telegram token url",
			requestBody:             `{"message":"opossum","chat_id":1234}`,
			botToken:                strPtr("bad\ntoken"),
			telegramShouldBeCalled:  false,
			expectedOutboundMessage: nil,
			responseCode:            http.StatusInternalServerError,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)
		botToken := strPtr("1")
		if testData.botToken != nil {
			botToken = testData.botToken
		}

		controller := Controller{
			Config: NewConfigSafe(strPtr("8080"), botToken, testData.defaultChatID, nil),
			HTTPClient: &MockHTTPClient{
				do: func(req *http.Request) (*http.Response, error) {
					if !testData.telegramShouldBeCalled {
						assert.Fail("do function should not be called")
					} else {
						m := messages.NewTelegramMessage(nil)
						err := json.NewDecoder(req.Body).Decode(&m)
						if err != nil {
							panic(fmt.Sprintf("Cannot decode outbound telegram message: %s", err))
						}

						assert.Equal(*testData.expectedOutboundMessage, m, "Outbound message is not as expected")
					}

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

		controller.SendMessage(w, req)

		assert.Equal(testData.responseCode, w.Code, fmt.Sprintf("Response code is not correct: %s", formatBody(w)))
	}
}

func formatBody(w *httptest.ResponseRecorder) string {
	bodyBytes, err := io.ReadAll(w.Body)
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

type errorReadCloser struct{}

func (errorReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errorReadCloser) Close() error {
	return nil
}

func TestTelegramControllerSendMessageCopyFailure(t *testing.T) {
	controller := Controller{
		Config: NewConfigSafe(strPtr("8080"), strPtr("1"), nil, nil),
		HTTPClient: &MockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"Content-Type":   []string{"application/json"},
						"Content-Length": []string{"10"},
					},
					Body: errorReadCloser{},
				}, nil
			},
		},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
		bytes.NewReader([]byte(`{"message":"hello","chat_id":1}`)))

	controller.SendMessage(w, req)

	if !strings.Contains(w.Body.String(), "Cannot copy a response.") {
		t.Fatalf("expected copy failure message in body, got %q", w.Body.String())
	}
}

func TestMessageConstructors(t *testing.T) {
	t.Run("new telegram message defaults", func(t *testing.T) {
		chatID := 42
		msg := NewTelegramMessage(&chatID)

		if msg.ChatID == nil || *msg.ChatID != chatID {
			t.Fatalf("expected chat id %d, got %+v", chatID, msg.ChatID)
		}
		if !msg.DisablePreview {
			t.Fatal("expected disable_web_page_preview to default to true")
		}
		if !msg.DisableNotification {
			t.Fatal("expected disable_notification to default to true")
		}
	})

	t.Run("new incoming message defaults", func(t *testing.T) {
		chatID := 77
		msg := NewMessage(&chatID)

		if msg.ChatID == nil || *msg.ChatID != chatID {
			t.Fatalf("expected chat id %d, got %+v", chatID, msg.ChatID)
		}
		if !msg.Silent {
			t.Fatal("expected silent to default to true")
		}
		if msg.Message != "" {
			t.Fatalf("expected empty message text by default, got %q", msg.Message)
		}
	})
}
