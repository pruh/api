package tests

import (
	"bytes"
	"fmt"

	"github.com/j-rooft/api/utils"
)

func NewConfig(port *string, botToken *string, defaultChatId *string, credsMap *map[string]string) (*utils.Configuration, error) {
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

	return utils.NewFromParams(port, botToken, defaultChatId, &authCreds)
}

func NewConfigSafe(port *string, botToken *string, defaultChatId *string, credsMap *map[string]string) *utils.Configuration {
	config, err := NewConfig(port, botToken, defaultChatId, credsMap)
	if err != nil {
		panic("Error is not nil")
	}

	return config
}
