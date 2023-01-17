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
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"apiVer": "23","type": 1,"omadacId": "someId"}}`,
			omadaErrorCode: 0,
			omadacId:       "someId",
		}, {
			description:    "error if upstream HTTP error",
			omadaUrl:       "https://omada.example.com",
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{}`,
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
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, NewStr(testData.omadaUrl), nil, nil),
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
		setCookie      http.Cookie
		loginToken     string
	}{
		{
			description:    "happy path",
			omadaUrl:       "https://omada.example.com",
			omadaLogin:     NewStr("username"),
			omadaPassword:  NewStr("password"),
			omadacId:       NewStr("omada_cid"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"apiVer": "23","type": 1,"token": "login_token"}}`,
			omadaErrorCode: 0,
			setCookie:      http.Cookie{Name: "rememberMe", Value: "deleteMe", Path: "/", MaxAge: 0},
			loginToken:     "login_token",
		},
		{
			description:    "upstream error",
			omadaUrl:       "https://omada.example.com",
			omadaLogin:     NewStr("username"),
			omadaPassword:  NewStr("password"),
			omadacId:       NewStr("omada_cid"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"apiVer": "23","type": 1,"token": "login_token"}}`,
			omadaErrorCode: 0,
			setCookie:      http.Cookie{Name: "rememberMe", Value: "deleteMe", Path: "/", MaxAge: 0},
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

				loginData := OmadaLoginData{}
				if err := json.NewDecoder(req.Body).Decode(&loginData); err != nil {
					panic(fmt.Sprintf("can not parse request body: %s", err))
				}

				defer req.Body.Close()

				assert.Equal(*testData.omadaLogin, loginData.Username, "omada username is not correct")
				assert.Equal(*testData.omadaPassword, loginData.Password, "omada password is not correct")

				if testData.responseCode != http.StatusOK {
					return nil, errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.Header().Add("Set-Cookie", testData.setCookie.String())
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}
				w.WriteHeader(testData.responseCode)

				return w.Result(), nil
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), testData.omadaLogin, testData.omadaPassword),
			&mockHttpClient)

		o, cookies, err := omadaApi.Login(testData.omadacId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, o.ErrorCode, "Error code is not correct")
		assert.Equal(testData.loginToken, *o.Result.Token, "Login token is not correct")

		assert.True(len(cookies) == 1, "cookies are not correct")
		assert.Equal(testData.setCookie.String(), cookies[0].String(), "cookies are not correct")
	}
}

func TestGetSites(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "site_name", "id": "site_id"}]}}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "site_name", "id": "site_id"}]}}`,
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

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		sitesResp, err := omadaApi.GetSites(testData.omadacId, testData.cookies, testData.loginToken)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, sitesResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*sitesResp.Result.Data), "sites response is not correct")
		assert.Equal(Data{Id: NewStr("site_id"), Name: NewStr("site_name")}, (*sitesResp.Result.Data)[0], "omada id parsing is not correct")
	}
}

func TestGetWlans(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		siteId         *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:         NewStr("site_id"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "wlan_name", "id": "wlan_id"}]}}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:         NewStr("site_id"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "wlan_name", "id": "wlan_id"}]}}`,
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

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		wlansResp, err := omadaApi.GetWlans(testData.omadacId, testData.cookies, testData.loginToken, testData.siteId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, wlansResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*wlansResp.Result.Data), "wlans resp is not correct")
		assert.Equal(Data{Id: NewStr("wlan_id"), Name: NewStr("wlan_name")}, (*wlansResp.Result.Data)[0], "omada id parsing is not correct")
	}
}

func TestGetSsids(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		siteId         *string
		wlanId         *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:         NewStr("site_id"),
			wlanId:         NewStr("wlan_id"),
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "ssid_name", "id": "ssid_id"}]}}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:         NewStr("site_id"),
			wlanId:         NewStr("wlan_id"),
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"errorCode": 0,"msg": "Success.","result": {"data": [{"name": "ssid_name", "id": "ssid_id"}]}}`,
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

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		wlansResp, err := omadaApi.GetSsids(testData.omadacId, testData.cookies, testData.loginToken, testData.siteId, testData.wlanId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, wlansResp.ErrorCode, "Error code is not correct")
		assert.Equal(1, len(*wlansResp.Result.Data), "ssids resp is not correct")
		assert.Equal(Data{Id: NewStr("ssid_id"), Name: NewStr("ssid_name")},
			(*wlansResp.Result.Data)[0], "omada id parsing is not correct")
	}
}

func TestUpdateSsid(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		siteId         *string
		wlanId         *string
		ssidId         *string
		ssidUpdateData *Data
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId: NewStr("site_id"),
			wlanId: NewStr("wlan_id"),
			ssidId: NewStr("ssid_id"),
			ssidUpdateData: &Data{Id: NewStr("ssid_id"), Name: NewStr("ssid_name"),
				WlanScheduleEnable: NewBool(true), Action: NewInt(0), ScheduleId: NewStr("schedule_id")},
			expectError:    false,
			responseCode:   http.StatusOK,
			responseBody:   `{"errorCode": 0,"msg": "Success."}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId: NewStr("site_id"),
			wlanId: NewStr("wlan_id"),
			ssidId: NewStr("ssid_id"),
			ssidUpdateData: &Data{Id: NewStr("ssid_id"), Name: NewStr("ssid_name"),
				WlanScheduleEnable: NewBool(true), Action: NewInt(0), ScheduleId: NewStr("schedule_id")},
			expectError:    true,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"errorCode": 0,"msg": "Success."}`,
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans/%s/ssids/%s",
					testData.omadaUrl, *testData.omadacId, *testData.siteId, *testData.wlanId, *testData.ssidId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				updateData := Data{}
				if err := json.NewDecoder(req.Body).Decode(&updateData); err != nil {
					panic(fmt.Sprintf("can not parse request body: %s", err))
				}

				defer req.Body.Close()

				assert.Equal(Data{
					Id:                 NewStr("ssid_id"),
					Name:               NewStr("ssid_name"),
					WlanScheduleEnable: NewBool(true),
					Action:             NewInt(0),
					ScheduleId:         NewStr("schedule_id"),
				}, updateData, "omada username is not correct")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		wlansResp, err := omadaApi.UpdateSsid(testData.omadacId, testData.cookies, testData.loginToken,
			testData.siteId, testData.wlanId, testData.ssidUpdateData)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, wlansResp.ErrorCode, "Error code is not correct")
	}
}

func TestGetTimeRanges(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		siteId         *string
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:       NewStr("site_id"),
			expectError:  false,
			responseCode: http.StatusOK,
			responseBody: `{
				"errorCode": 0,"msg": "Success.","result": {
					"data": [{
						"name": "tr_name", "id": "tr_id", "daymode": 1, "timelist": [{
							"startTimeH": 1, "startTimeM": 22, "endTimeH": 2, "endTimeM": 55
						}]
					}]
				}
			}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId:       NewStr("site_id"),
			expectError:  true,
			responseCode: http.StatusInternalServerError,
			responseBody: `{
				"errorCode": 0,"msg": "Success.","result": {
					"data": [{
						"name": "tr_name", "id": "tr_id", "daymode": 1, "timelist": [{
							"startTimeH": 1, "startTimeM": 22, "endTimeH": 2, "endTimeM": 55
						}]
					}]
				}
			}`,
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/profiles/timeranges",
					testData.omadaUrl, *testData.omadacId, *testData.siteId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		trResp, err := omadaApi.GetTimeRanges(testData.omadacId, testData.cookies, testData.loginToken,
			testData.siteId)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, trResp.ErrorCode, "Error code is not correct")
		assert.Equal(Data{
			Id:      NewStr("tr_id"),
			Name:    NewStr("tr_name"),
			DayMode: NewInt(1),
			TimeList: &[]TimeList{
				{
					StartTimeH: NewInt(1),
					StartTimeM: NewInt(22),
					EndTimeH:   NewInt(2),
					EndTimeM:   NewInt(55),
				},
			},
		}, (*trResp.Result.Data)[0], "Error code is not correct")
	}
}

func TestCreateTimeRange(t *testing.T) {
	testsData := []struct {
		description    string
		omadaUrl       string
		omadacId       *string
		loginToken     *string
		cookies        []*http.Cookie
		siteId         *string
		timeRangeData  *Data
		expectError    bool
		responseCode   int
		responseBody   string
		omadaErrorCode int
	}{
		{
			description: "happy path",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId: NewStr("site_id"),
			timeRangeData: &Data{
				Name:    NewStr("Night and Day"),
				DayMode: NewInt(0),
				DayMon:  NewBool(true),
				DayTue:  NewBool(true),
				DayWed:  NewBool(true),
				DayThu:  NewBool(true),
				DayFri:  NewBool(true),
				DaySat:  NewBool(true),
				DaySun:  NewBool(true),
				TimeList: &[]TimeList{
					{
						DayType:    NewInt(0),
						StartTimeH: NewInt(0),
						StartTimeM: NewInt(0),
						EndTimeH:   NewInt(24),
						EndTimeM:   NewInt(0),
					},
				}},
			expectError:  false,
			responseCode: http.StatusOK,
			responseBody: `{
				"errorCode": 0,"msg": "Success.","result": {
					"data": [{
						"name": "tr_name", "id": "tr_id", "daymode": 1, "timelist": [{
							"startTimeH": 1, "startTimeM": 22, "endTimeH": 2, "endTimeM": 55
						}]
					}]
				}
			}`,
			omadaErrorCode: 0,
		},
		{
			description: "upstream error",
			omadaUrl:    "https://omada.example.com",
			omadacId:    NewStr("omada_cid"),
			loginToken:  NewStr("login_token"),
			cookies: []*http.Cookie{
				{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			},
			siteId: NewStr("site_id"),
			timeRangeData: &Data{
				Name:    NewStr("Night and Day"),
				DayMode: NewInt(0),
				DayMon:  NewBool(true),
				DayTue:  NewBool(true),
				DayWed:  NewBool(true),
				DayThu:  NewBool(true),
				DayFri:  NewBool(true),
				DaySat:  NewBool(true),
				DaySun:  NewBool(true),
				TimeList: &[]TimeList{
					{
						DayType:    NewInt(0),
						StartTimeH: NewInt(0),
						StartTimeM: NewInt(0),
						EndTimeH:   NewInt(24),
						EndTimeM:   NewInt(0),
					},
				}},
			expectError:  true,
			responseCode: http.StatusInternalServerError,
			responseBody: `{
				"errorCode": 0,"msg": "Success.","result": {
					"data": [{
						"name": "tr_name", "id": "tr_id", "daymode": 1, "timelist": [{
							"startTimeH": 1, "startTimeM": 22, "endTimeH": 2, "endTimeM": 55
						}]
					}]
				}
			}`,
			omadaErrorCode: 0,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		mockHttpClient := MockHTTPClient{
			MockDo: func(req *http.Request) (*http.Response, error) {
				assert.Equal(fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/profiles/timeranges",
					testData.omadaUrl, *testData.omadacId, *testData.siteId),
					req.URL.String(), "Omada request url is not correct")

				assert.Equal(testData.cookies, req.Cookies(), "cookie is missing")
				assert.Equal(*testData.loginToken, req.Header.Get("Csrf-token"), "Login token is missing")

				var respErr error
				if testData.responseCode != http.StatusOK {
					respErr = errors.New("test error")
				}

				w := httptest.NewRecorder()
				w.WriteHeader(testData.responseCode)
				if _, err := w.WriteString(testData.responseBody); err != nil {
					panic(fmt.Sprintf("Error writing body: %s", err))
				}

				return w.Result(), respErr
			},
		}

		omadaApi := NewOmadaApi(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.omadaUrl), nil, nil),
			&mockHttpClient)

		trResp, err := omadaApi.CreateTimeRange(testData.omadacId, testData.cookies, testData.loginToken,
			testData.siteId, testData.timeRangeData)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(testData.omadaErrorCode, trResp.ErrorCode, "Error code is not correct")
		assert.Equal(Data{
			Id:      NewStr("tr_id"),
			Name:    NewStr("tr_name"),
			DayMode: NewInt(1),
			TimeList: &[]TimeList{
				{
					StartTimeH: NewInt(1),
					StartTimeM: NewInt(22),
					EndTimeH:   NewInt(2),
					EndTimeM:   NewInt(55),
				},
			},
		}, (*trResp.Result.Data)[0], "Error code is not correct")
	}
}
