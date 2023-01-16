package networks_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/pruh/api/config/tests"
	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
	"github.com/stretchr/testify/assert"
)

func TestGetControllerId(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
		omadacId       string
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"apiVer\": \"23\",\"type\": 1,\"omadacId\": \"someId\"}}",
			omadaErrorCode: 0,
			omadacId:       "someId",
		}, {
			description:    "error if upstream HTTP error",
			omadaUrl:       "https://omada.example.com",
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   "{}",
			omadaErrorCode: 0,
			omadacId:       "someId",
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(testData.omadaUrl+"/api/info", req.URL.String(), "Omada request url is not correct")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				w.WriteString(testData.responseBody)

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil, StrPtr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		omadaControllerId, err := omadaApi.GetControllerId()
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, omadaControllerId.ErrorCode, "Error code is not correct")
		assert.Equal(testData.omadacId, *omadaControllerId.Result.OmadacId, "Omada id parsing is not correct")
	}
}

func TestLogin(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadaLogin     *string
		omadaPassword  *string
		omadacId       *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
		loginToken     string
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			omadaLogin:     StrPtr("username"),
			omadaPassword:  StrPtr("password"),
			omadacId:       StrPtr("omada_cid"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"apiVer\": \"23\",\"type\": 1,\"token\": \"login_token\"}}",
			omadaErrorCode: 0,
			loginToken:     "login_token",
		},
		{
			description:    "upstream error",
			omadaUrl:       "https://omada.example.com",
			omadaLogin:     StrPtr("username"),
			omadaPassword:  StrPtr("password"),
			omadacId:       StrPtr("omada_cid"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"apiVer\": \"23\",\"type\": 1,\"token\": \"login_token\"}}",
			omadaErrorCode: 0,
			loginToken:     "login_token",
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/login", testData.omadaUrl, *testData.omadacId),
					req.URL.String(), "Omada request url is not correct")

				loginData := LoginData{}
				if err := json.NewDecoder(req.Body).Decode(&loginData); err != nil {
					panic(fmt.Sprintf("can not parse request body: %s", err))
				}

				defer req.Body.Close()

				assert.Equal(*testData.omadaLogin, loginData.Username, "omada username is not correct")
				assert.Equal(*testData.omadaPassword, loginData.Password, "omada password is not correct")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				w.WriteString(testData.responseBody)

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil,
				StrPtr(testData.omadaUrl), testData.omadaLogin, testData.omadaPassword),
			&mockHttpClient)

		omadaControllerId, err := omadaApi.Login(testData.omadacId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, omadaControllerId.ErrorCode, "Error code is not correct")
		assert.Equal(testData.loginToken, *omadaControllerId.Result.Token, "Login token is not correct")
	}
}

func TestGetSites(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"site_name\", \"id\": \"site_id\"}]}}",
			omadaErrorCode: 0,
		},
		{
			description:    "upstream error",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"site_name\", \"id\": \"site_id\"}]}}",
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites?currentPageSize=1&currentPage=1",
					testData.omadaUrl, *testData.omadacId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				w.WriteString(testData.responseBody)

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil,
				StrPtr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		sitesResp, err := omadaApi.GetSites(testData.omadacId, testData.loginToken)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, sitesResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*sitesResp.Result.Data), "sites response is not correct")
		assert.Equal(Data{Id: StrPtr("site_id"), Name: StrPtr("site_name")}, (*sitesResp.Result.Data)[0], "omada id parsing is not correct")
	}
}

func TestGetWlans(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		siteId         *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			siteId:         StrPtr("site_id"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"wlan_name\", \"id\": \"wlan_id\"}]}}",
			omadaErrorCode: 0,
		},
		{
			description:    "upstream error",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			siteId:         StrPtr("site_id"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"wlan_name\", \"id\": \"wlan_id\"}]}}",
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans",
					testData.omadaUrl, *testData.omadacId, *testData.siteId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				w.WriteString(testData.responseBody)

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil,
				StrPtr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		wlansResp, err := omadaApi.GetWlans(testData.omadacId, testData.loginToken, testData.siteId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, wlansResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*wlansResp.Result.Data), "wlans resp is not correct")
		assert.Equal(Data{Id: StrPtr("wlan_id"), Name: StrPtr("wlan_name")}, (*wlansResp.Result.Data)[0], "omada id parsing is not correct")
	}
}

func TestGetSsids(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		siteId         *string
		wlanId         *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			siteId:         StrPtr("site_id"),
			wlanId:         StrPtr("wlan_id"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"ssid_name\", \"id\": \"ssid_id\"}]}}",
			omadaErrorCode: 0,
		},
		{
			description:    "upstream error",
			omadaUrl:       "https://omada.example.com",
			omadacId:       StrPtr("omada_cid"),
			loginToken:     StrPtr("login_token"),
			siteId:         StrPtr("site_id"),
			wlanId:         StrPtr("wlan_id"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   "{\"errorCode\": 0,\"msg\": \"Success.\",\"result\": {\"data\": [{\"name\": \"ssid_name\", \"id\": \"ssid_id\"}]}}",
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans/%s/ssids",
					testData.omadaUrl, *testData.omadacId, *testData.siteId, *testData.wlanId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				w.WriteString(testData.responseBody)

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil,
				StrPtr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		wlansResp, err := omadaApi.GetSsids(testData.omadacId, testData.loginToken, testData.siteId, testData.wlanId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, wlansResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*wlansResp.Result.Data), "ssids resp is not correct")
		assert.Equal(Data{Id: StrPtr("ssid_id"), Name: StrPtr("ssid_name")},
			(*wlansResp.Result.Data)[0], "omada id parsing is not correct")
	}
}
