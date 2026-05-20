package config_test

import (
	"testing"

	"github.com/pruh/api/config"
	. "github.com/pruh/api/config/tests"
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

		conf, err := NewConfig(testData.port, testData.botToken, testData.defaultChatID,
			testData.credsMap)

		if !testData.expectError && err != nil {
			assert.Fail("Config load should not return error: " + err.Error())
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

func TestNewFromEnv(t *testing.T) {
	t.Run("missing required variables", func(t *testing.T) {
		t.Setenv("PORT", "")
		t.Setenv("TELEGRAM_BOT_TOKEN", "")
		t.Setenv("TELEGRAM_DEFAULT_CHAT_ID", "")
		t.Setenv("API_V1_CREDS", "")

		cfg, err := config.NewFromEnv()
		if err == nil {
			t.Fatalf("expected error, got config %+v", cfg)
		}
	})

	t.Run("invalid default chat id", func(t *testing.T) {
		t.Setenv("PORT", "8080")
		t.Setenv("TELEGRAM_BOT_TOKEN", "token")
		t.Setenv("TELEGRAM_DEFAULT_CHAT_ID", "not-a-number")
		t.Setenv("API_V1_CREDS", "")

		cfg, err := config.NewFromEnv()
		if err == nil {
			t.Fatalf("expected error, got config %+v", cfg)
		}
	})

	t.Run("invalid credentials json", func(t *testing.T) {
		t.Setenv("PORT", "8080")
		t.Setenv("TELEGRAM_BOT_TOKEN", "token")
		t.Setenv("TELEGRAM_DEFAULT_CHAT_ID", "")
		t.Setenv("API_V1_CREDS", "not-json")

		cfg, err := config.NewFromEnv()
		if err == nil {
			t.Fatalf("expected error, got config %+v", cfg)
		}
	})

	t.Run("valid optional values", func(t *testing.T) {
		t.Setenv("PORT", "8080")
		t.Setenv("TELEGRAM_BOT_TOKEN", "token")
		t.Setenv("TELEGRAM_DEFAULT_CHAT_ID", "1234")
		t.Setenv("API_V1_CREDS", `{"alice":"secret"}`)

		cfg, err := config.NewFromEnv()
		if err != nil {
			t.Fatalf("did not expect error: %v", err)
		}
		if cfg.DefaultChatID == nil || *cfg.DefaultChatID != 1234 {
			t.Fatalf("expected default chat id 1234, got %+v", cfg.DefaultChatID)
		}
		if cfg.APIV1Credentials == nil {
			t.Fatal("expected credentials to be set")
		}
		if got := (*cfg.APIV1Credentials)["alice"]; got != "secret" {
			t.Fatalf("expected alice credential to be secret, got %q", got)
		}
	})
}
