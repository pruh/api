package utils

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

// Configuration contrains configuration parameters.
type Configuration struct {
	Port             *string
	TelegramBoToken  *string
	DefaultChatID    *int
	APIV1Credentials *map[string]string
}

// NewFromEnv creates new configuration from environment variables.
func NewFromEnv() (*Configuration, error) {
	port := ptr(os.Getenv("PORT"))
	botToken := ptr(os.Getenv("TELEGRAM_BOT_TOKEN"))
	chatID := ptrOrNil(os.LookupEnv("TELEGRAM_DEFAULT_CHAT_ID"))
	apiCreds := ptrOrNil(os.LookupEnv("API_V1_CREDS"))

	return NewFromParams(port, botToken, chatID, apiCreds)
}

func ptr(str string) *string {
	return &str
}

func ptrOrNil(str string, valueSet bool) *string {
	if valueSet {
		return &str
	}

	return nil
}

// NewFromEnv creates new configuration from arguments.
func NewFromParams(port *string, boToken *string, defaultChatID *string, apiV1Credentials *string) (*Configuration, error) {
	var conf Configuration
	if port == nil || *port == "" {
		return nil, errors.New("port should not be empty")
	}
	conf.Port = port

	if boToken == nil || *boToken == "" {
		return nil, errors.New("telegram bot token should not be empty")
	}
	conf.TelegramBoToken = boToken

	if defaultChatID != nil && *defaultChatID != "" {
		if chatID, err := strconv.Atoi(*defaultChatID); err == nil {
			conf.DefaultChatID = &chatID
		} else {
			return nil, err
		}
	}

	if apiV1Credentials != nil && *apiV1Credentials != "" {
		conf.APIV1Credentials = &map[string]string{}
		err := json.Unmarshal([]byte(*apiV1Credentials), conf.APIV1Credentials)
		if err != nil {
			return nil, err
		}
	}

	return &conf, nil
}
