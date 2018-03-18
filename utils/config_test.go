package utils_test

import (
	"testing"

	. "github.com/pruh/api/tests"
	"github.com/stretchr/testify/assert"
)

func TestConfiguraionLoading(t *testing.T) {
	testsData := []struct {
		description   string
		port          *string
		botToken      *string
		defaultChatID *string
		credsMap      *map[string]string
		expectError   bool
	}{
		{
			description:   "normal config",
			port:          ptr("1234"),
			botToken:      ptr("botToken"),
			defaultChatID: ptr("1234"),
			credsMap: &map[string]string{
				"bob":  "dylan",
				"jack": "sparrow",
			},
			expectError: false,
		},
		{
			description:   "nil port",
			port:          nil,
			botToken:      ptr("botToken"),
			defaultChatID: ptr("1234"),
			credsMap: &map[string]string{
				"bob":  "dylan",
				"jack": "sparrow",
			},
			expectError: true,
		},
		{
			description:   "empty port",
			port:          ptr(""),
			botToken:      ptr("botToken"),
			defaultChatID: ptr("1234"),
			credsMap: &map[string]string{
				"bob":  "dylan",
				"jack": "sparrow",
			},
			expectError: true,
		},
		{
			description:   "nil bot token",
			port:          ptr("1234"),
			botToken:      nil,
			defaultChatID: ptr("1234"),
			credsMap: &map[string]string{
				"bob":  "dylan",
				"jack": "sparrow",
			},
			expectError: true,
		},
		{
			description:   "empty bot token",
			port:          ptr("1234"),
			botToken:      ptr(""),
			defaultChatID: ptr("1234"),
			credsMap: &map[string]string{
				"bob":  "dylan",
				"jack": "sparrow",
			},
			expectError: true,
		},
		{
			description:   "nil default chat id and credentials",
			port:          ptr("1234"),
			botToken:      ptr("botToken"),
			defaultChatID: nil,
			credsMap:      nil,
			expectError:   false,
		},
		{
			description:   "empty default chat id and credentials",
			port:          ptr("1234"),
			botToken:      ptr("botToken"),
			defaultChatID: ptr(""),
			credsMap:      &map[string]string{},
			expectError:   false,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("testing %+v", testData.description)

		conf, err := NewConfig(testData.port, testData.botToken, testData.defaultChatID, testData.credsMap)

		if !testData.expectError && err != nil {
			assert.Fail("Config load should not return error: %s", err.Error)
			continue
		}

		if testData.expectError && err == nil {
			assert.Fail("Config load should return error")
			continue
		}

		if testData.expectError {
			continue
		}

		assert.Equal(testData.port, conf.Port, "Port is not correct")
		assert.Equal(testData.botToken, conf.TelegramBoToken, "Bot token is not correct")
		assert.Equal(testData.credsMap, conf.APIV1Credentials, "Credentials is not correct")
	}
}

func ptr(str string) *string {
	return &str
}
