package tests

import (
	"bytes"
	"fmt"

	"github.com/pruh/api/config"
)

// NewConfig creates new config or returns error.
func NewConfig(port *string, botToken *string, defaultChatID *string, credsMap *map[string]string) (*config.Configuration, error) {
	var buffer bytes.Buffer
	var authCreds string
	if credsMap != nil {
		buffer.WriteString("{")
		for k, v := range *credsMap {
			if buffer.Len() > 1 {
				buffer.WriteString(",")
			}
			buffer.WriteString(fmt.Sprintf(`"%s":"%s"`, k, v))
		}
		buffer.WriteString("}")
		authCreds = buffer.String()
	}

	return config.NewFromParams(port, botToken, defaultChatID, &authCreds)
}

// NewConfigSafe creates new config or calls panic if config can not be created.
func NewConfigSafe(port *string, botToken *string, defaultChatID *string, credsMap *map[string]string) *config.Configuration {
	config, err := NewConfig(port, botToken, defaultChatID, credsMap)
	if err != nil {
		panic("Error is not nil")
	}

	return config
}