package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pruh/api/utils"
	"github.com/stretchr/testify/assert"

	. "github.com/pruh/api/middleware"
	. "github.com/pruh/api/tests"
)

func TestBasicAuth(t *testing.T) {
	testsData := []struct {
		description  string
		user         string
		password     string
		config       *utils.Configuration
		requestBody  io.Reader
		responseCode int
	}{
		{
			description: "happy path",
			user:        "papa",
			password:    "castoro",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"papa": "castoro",
			}),
			requestBody:  nil,
			responseCode: http.StatusOK,
		},
		{
			description: "non-empty body",
			user:        "papa",
			password:    "castoro",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"papa": "castoro",
			}),
			requestBody:  bytes.NewReader([]byte(`{"test":"test"}`)),
			responseCode: http.StatusOK,
		},
		{
			description: "multiple credential",
			user:        "papa",
			password:    "castoro",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"mama": "castoro",
				"papa": "castoro",
			}),
			requestBody:  nil,
			responseCode: http.StatusOK,
		},
		{
			description: "wrong username",
			user:        "mama",
			password:    "castoro",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"papa": "castoro",
			}),
			requestBody:  nil,
			responseCode: http.StatusUnauthorized,
		},
		{
			description: "wrong password",
			user:        "papa",
			password:    "castoro2",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"papa": "castoro",
			}),
			requestBody:  nil,
			responseCode: http.StatusUnauthorized,
		},
		{
			description:  "empty credentials",
			user:         "",
			password:     "",
			config:       NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{}),
			requestBody:  nil,
			responseCode: http.StatusOK,
		},
		{
			description:  "nil credentials",
			user:         "",
			password:     "",
			config:       NewConfigSafe(ptr("8080"), ptr("1"), nil, nil),
			requestBody:  nil,
			responseCode: http.StatusOK,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("testing %+v", testData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", testData.requestBody)
		if testData.user != "" || testData.password != "" {
			req.SetBasicAuth(testData.user, testData.password)
		}

		AuthMiddleware(w, req, func(w http.ResponseWriter, r *http.Request) {
			if testData.responseCode != http.StatusOK {
				assert.Fail("next handler should not be called when testing %s", testData.description)
			}
		}, testData.config)

		assert.Equalf(testData.responseCode, w.Code, "response code is not correct for %s test", testData.description)
	}
}

func ptr(str string) *string {
	return &str
}
