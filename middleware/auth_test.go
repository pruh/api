package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j-rooft/api/utils"
	"github.com/stretchr/testify/assert"

	. "github.com/j-rooft/api/middleware"
	. "github.com/j-rooft/api/tests"
)

func TestBasicAuth(t *testing.T) {
	testsData := []struct {
		description  string
		user         string
		password     string
		config       *utils.Configuration
		remoteIP     string
		xFwdHeader   string
		xRealIP      string
		requestBody  io.Reader
		responseCode int
	}{
		// {
		// 	description: "happy path",
		// 	user:        "papa",
		// 	password:    "castoro",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description: "non-empty body",
		// 	user:        "papa",
		// 	password:    "castoro",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	requestBody:  bytes.NewReader([]byte(`{"test":"test"}`)),
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description: "multiple credential",
		// 	user:        "papa",
		// 	password:    "castoro",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"mama": "castoro",
		// 		"papa": "castoro",
		// 	}),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description: "wrong username",
		// 	user:        "mama",
		// 	password:    "castoro",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusUnauthorized,
		// },
		// {
		// 	description: "wrong password",
		// 	user:        "papa",
		// 	password:    "castoro2",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusUnauthorized,
		// },
		// {
		// 	description:  "empty credentials",
		// 	user:         "",
		// 	password:     "",
		// 	config:       NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{}),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description:  "nil credentials",
		// 	user:         "",
		// 	password:     "",
		// 	config:       NewConfigSafe(ptr("8080"), ptr("1"), nil, nil),
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		{
			description: "local network request",
			user:        "",
			password:    "",
			config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
				"papa": "castoro",
			}),
			remoteIP:     "192.168.0.2:8080",
			requestBody:  nil,
			responseCode: http.StatusOK,
		},
		// {
		// 	description: "local network request for X-Forwarded-For local",
		// 	user:        "",
		// 	password:    "",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	remoteIP:     "8.8.4.4:8080",
		// 	xFwdHeader:   "192.168.1.2, 8.8.4.4, 10.8.0.1",
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description: "local network request for X-Forwarded-For remote",
		// 	user:        "",
		// 	password:    "",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	xFwdHeader:   "10.8.0.1, 192.168.1.2, 8.8.4.4",
		// 	requestBody:  nil,
		// 	responseCode: http.StatusUnauthorized,
		// },
		// {
		// 	description: "local network request for X-Real-IP local",
		// 	user:        "",
		// 	password:    "",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	remoteIP:     "8.8.4.4",
		// 	xRealIP:      "10.8.0.1",
		// 	requestBody:  nil,
		// 	responseCode: http.StatusOK,
		// },
		// {
		// 	description: "local network request for X-Real-IP remote",
		// 	user:        "",
		// 	password:    "",
		// 	config: NewConfigSafe(ptr("8080"), ptr("1"), nil, &map[string]string{
		// 		"papa": "castoro",
		// 	}),
		// 	xRealIP:      "8.8.4.4",
		// 	requestBody:  nil,
		// 	responseCode: http.StatusUnauthorized,
		// },
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("testing %+v", testData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", testData.requestBody)
		if testData.user != "" || testData.password != "" {
			req.SetBasicAuth(testData.user, testData.password)
		}
		if testData.remoteIP != "" {
			req.RemoteAddr = testData.remoteIP
		} else {
			req.RemoteAddr = "8.8.8.8"
		}
		if testData.xFwdHeader != "" {
			req.Header.Set("X-Forwarded-For", testData.xFwdHeader)
		}
		if testData.xRealIP != "" {
			req.Header.Set("X-Real-IP", testData.xRealIP)
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
