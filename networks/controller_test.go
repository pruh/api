package networks_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	. "github.com/pruh/api/config/tests"
	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
)

func TestGetWifi(t *testing.T) {
	testsData := []struct {
		description string
		requestUrl  string
		ssidParam   *string
		cidError    bool
		loginError  bool
		siteError   bool
		wlanError   bool
		ssidError   bool

		responseCode int
	}{
		{
			description:  "ControllerId happy path",
			requestUrl:   "https://omada.example.com/networks/ssid",
			ssidParam:    NewStr("my_ssid"),
			responseCode: http.StatusOK,
		},
		{
			description:  "ssid missing in the request params",
			requestUrl:   "https://omada.example.com",
			responseCode: http.StatusBadRequest,
		},
		{
			description: "omada controller id response error",
			requestUrl:  "https://omada.example.com",
			ssidParam:   NewStr("my_ssid"),
			cidError:    true,

			responseCode: http.StatusBadGateway,
		},
		{
			description: "omada login query error",
			requestUrl:  "https://omada.example.com",
			ssidParam:   NewStr("my_ssid"),
			loginError:  true,

			responseCode: http.StatusBadGateway,
		},
		{
			description: "omada site response error",
			requestUrl:  "https://omada.example.com",
			ssidParam:   NewStr("my_ssid"),
			siteError:   true,

			responseCode: http.StatusBadGateway,
		},
		{
			description: "omada wlan response error",
			requestUrl:  "https://omada.example.com",
			ssidParam:   NewStr("my_ssid"),
			wlanError:   true,

			responseCode: http.StatusBadGateway,
		},
		{
			description: "omada ssid response error",
			requestUrl:  "https://omada.example.com",
			ssidParam:   NewStr("my_ssid"),
			ssidError:   true,

			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					if testData.cidError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("oc_id"),
						},
					}

					return resp, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					if testData.loginError {
						return nil, nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}
					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string) (*OmadaResponse, error) {
					if testData.siteError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("site_id"),
									Name: NewStr("site_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string, siteId *string) (*OmadaResponse, error) {
					if testData.wlanError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					if testData.ssidError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:                 NewStr("ssid_id"),
									Name:               NewStr("my_ssid"),
									WlanScheduleEnable: NewBool(false),
									RateLimit: &RateLimit{
										UpLimitEnable:   NewBool(false),
										DownLimitEnable: NewBool(false),
									},
								},
							},
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

		// setting mux vars for testing
		if testData.ssidParam != nil {
			vars := map[string]string{
				"ssid": *testData.ssidParam,
			}
			req = mux.SetURLVars(req, vars)
		}

		controller.GetWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_ControllerId(t *testing.T) {
	testsData := []struct {
		description         string
		requestUrl          string
		requestData         *string
		ssidParam           *string
		omadaResponseError  bool
		omadaControllerId   *string
		loginToken          *string
		responseCode        int
		responseSsidUpdated *bool
		responseRadioOn     *bool
	}{
		{
			description:         "ControllerId happy path",
			requestUrl:          "https://omada.example.com/networks/ssid",
			requestData:         NewStr(`{"radioOn":false}`),
			ssidParam:           NewStr("my_ssid"),
			omadaControllerId:   NewStr("c_id"),
			loginToken:          NewStr("login_token"),
			responseCode:        http.StatusOK,
			responseSsidUpdated: NewBool(true),
			responseRadioOn:     NewBool(false),
		},
		{
			description:  "ssid missing in the request params",
			requestUrl:   "https://omada.example.com",
			requestData:  NewStr(`{"radioOn":false}`),
			responseCode: http.StatusBadRequest,
		},
		{
			description:        "omada controller id response error",
			requestUrl:         "https://omada.example.com",
			requestData:        NewStr(`{"radioOn":false}`),
			ssidParam:          NewStr("my_ssid"),
			omadaResponseError: true,
			omadaControllerId:  NewStr("c_id"),
			responseCode:       http.StatusBadGateway,
		},
		{
			description:       "request json malformed",
			requestUrl:        "https://omada.example.com",
			requestData:       NewStr(`aaa`),
			ssidParam:         NewStr("my_ssid"),
			omadaControllerId: nil,
			responseCode:      http.StatusBadRequest,
		},
		{
			description:         "request empty json",
			requestUrl:          "https://omada.example.com",
			requestData:         NewStr(`{}`),
			ssidParam:           NewStr("my_ssid"),
			omadaControllerId:   nil,
			responseCode:        http.StatusOK,
			responseSsidUpdated: NewBool(false),
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil,
				NewStr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: testData.omadaControllerId,
						},
					}

					return resp, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: testData.loginToken,
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}
					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("site_id"),
									Name: NewStr("site_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:                 NewStr("ssid_id"),
									Name:               NewStr("my_ssid"),
									WlanScheduleEnable: NewBool(false),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "http://example.com/foo",
			bytes.NewBuffer([]byte(*testData.requestData)))

		// setting mux vars for testing
		if testData.ssidParam != nil {
			vars := map[string]string{
				"ssid": *testData.ssidParam,
			}
			req = mux.SetURLVars(req, vars)
		}

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.Equal(*testData.responseSsidUpdated, *netsResponse.Updated, "Response updated flag is incorrect")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			if testData.responseRadioOn != nil {
				assert.Equal(*testData.responseRadioOn, *netsResponse.RadioOn, "Response success body is incorrect")
			}
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_Login(t *testing.T) {
	testsData := []struct {
		description        string
		requestUrl         string
		omadaResponseError bool
		loginToken         *string
		responseCode       int
	}{
		{
			description:  "Login happy path",
			requestUrl:   "https://omada.example.com/networks/ssid",
			loginToken:   NewStr("login_token"),
			responseCode: http.StatusOK,
		},
		{
			description:        "omada Login response error",
			requestUrl:         "https://omada.example.com",
			loginToken:         NewStr("login_token"),
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:  "LoginToken is missing in omada response",
			requestUrl:   "https://omada.example.com",
			loginToken:   nil,
			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, NewStr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					if testData.omadaResponseError {
						return nil, nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: testData.loginToken,
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:                 NewStr("ssid_id"),
									Name:               NewStr("my_ssid"),
									WlanScheduleEnable: NewBool(false),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_GetSites(t *testing.T) {
	testsData := []struct {
		description        string
		requestUrl         string
		omadaResponseError bool
		includeSites       bool
		responseCode       int
	}{
		{
			description:  "GetSites happy path",
			requestUrl:   "https://omada.example.com/networks/ssid",
			includeSites: true,
			responseCode: http.StatusOK,
		},
		{
			description:        "omada GetSites response error",
			requestUrl:         "https://omada.example.com",
			includeSites:       true,
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:  "site is missing in omada response",
			requestUrl:   "https://omada.example.com",
			includeSites: false,
			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, NewStr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					var res *Result
					if testData.includeSites {
						res = &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						}
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result:    res,
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:                 NewStr("ssid_id"),
									Name:               NewStr("my_ssid"),
									WlanScheduleEnable: NewBool(false),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_GetWlans(t *testing.T) {
	testsData := []struct {
		description        string
		omadaResponseError bool
		includeWlans       bool
		responseCode       int
	}{
		{
			description:  "GetWlans happy path",
			includeWlans: true,
			responseCode: http.StatusOK,
		},
		{
			description:        "omada GetSites response error",
			includeWlans:       true,
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:  "wlan is missing in omada response",
			includeWlans: false,
			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, nil, nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string, siteId *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					var res *Result
					if testData.includeWlans {
						res = &Result{
							Data: &[]Data{{Id: NewStr("wlan_id"), Name: NewStr("wlan_name")}},
						}
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result:    res,
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:                 NewStr("ssid_id"),
									Name:               NewStr("my_ssid"),
									WlanScheduleEnable: NewBool(false),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_GetSsids(t *testing.T) {
	testsData := []struct {
		description        string
		omadaResponseError bool
		includeSsids       bool
		responseCode       int
	}{
		{
			description:  "GetSsids happy path",
			includeSsids: true,
			responseCode: http.StatusOK,
		},
		{
			description:        "omada GetSites response error",
			includeSsids:       true,
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:  "wlan is missing in omada response",
			includeSsids: false,
			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, nil, nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					var res *Result
					if testData.includeSsids {
						res = &Result{
							Data: &[]Data{{
								Id:                 NewStr("ssid_id"),
								Name:               NewStr("my_ssid"),
								WlanScheduleEnable: NewBool(false),
							}},
						}
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result:    res,
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_UpdateSsid(t *testing.T) {
	testsData := []struct {
		description        string
		omadaResponseError bool
		includeSsids       bool
		responseCode       int
	}{
		{
			description:  "UpdateSsid happy path",
			includeSsids: true,
			responseCode: http.StatusOK,
		},
		{
			description:        "omada UpdateSsid response error",
			includeSsids:       true,
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, nil, nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{
								Id:                 NewStr("ssid_id"),
								Name:               NewStr("my_ssid"),
								WlanScheduleEnable: NewBool(false),
							}},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {

					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_GetTimeRanges(t *testing.T) {
	testsData := []struct {
		description        string
		omadaResponseError bool
		responseCode       int
	}{
		{
			description:  "GetTimeRanges happy path",
			responseCode: http.StatusOK,
		},
		{
			description:        "omada get time range response error",
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, nil, nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{
								Id:                 NewStr("ssid_id"),
								Name:               NewStr("my_ssid"),
								WlanScheduleEnable: NewBool(false),
							}},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {

					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("profile_id"),
							Data:      &[]Data{{Id: NewStr("time_range_id"), Name: NewStr("time_range_id")}},
						},
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*netsResponse.Updated, "Response success body missing updated flag")
			assert.True(netsResponse.Ssid != nil, "Response success body is incorrect")
			assert.True(netsResponse.RadioOn != nil, "Response success body is incorrect")
		} else {
			assert.True(len(*netsResponse.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func TestUpdateWifis_CreateTimeRanges(t *testing.T) {
	testsData := []struct {
		description            string
		omadaResponseError     bool
		returnCorrectTimeRange bool
		returnWrongTimeRange   bool
		responseCode           int
	}{
		{
			description:            "CreateTimeRanges happy path",
			returnCorrectTimeRange: true,
			responseCode:           http.StatusOK,
		},
		{
			description:        "omada get time range response error",
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:            "omada correct time range exists",
			returnCorrectTimeRange: true,
			responseCode:           http.StatusOK,
		},
		{
			description:          "omada wrong time range exists",
			returnWrongTimeRange: true,
			responseCode:         http.StatusOK,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, nil, nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							OmadacId: NewStr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					cookies := []*http.Cookie{
						{
							Name:  "cookie_name",
							Value: "cookie_value",
						},
					}

					return resp, cookies, nil
				},
				MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie,
					loginToken *string, siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("wlan_id"),
									Name: NewStr("wlan_name"),
								},
							},
						},
					}

					return resp, nil
				},
				MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{
								Id:                 NewStr("ssid_id"),
								Name:               NewStr("my_ssid"),
								WlanScheduleEnable: NewBool(false),
							}},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string) (*OmadaResponse, error) {

					res := &Result{}
					if testData.returnCorrectTimeRange || testData.returnWrongTimeRange {
						var endTime *int
						if testData.returnWrongTimeRange {
							endTime = NewInt(2)
						} else {
							endTime = NewInt(24)
						}
						res = &Result{
							ProfileId: NewStr("profile_id"),
							Data: &[]Data{
								{
									Id:      NewStr("time_range_id"),
									Name:    NewStr("time_range_id"),
									DayMode: NewInt(0),
									TimeList: &[]TimeList{
										{
											StartTimeH: NewInt(0),
											StartTimeM: NewInt(0),
											EndTimeH:   endTime,
											EndTimeM:   NewInt(0),
										},
									},
								},
							},
						}
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result:    res,
					}

					return resp, nil
				},
				MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
					siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					if testData.returnCorrectTimeRange {
						assert.Fail("should not be called for pre-defined time range")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							ProfileId: NewStr("time_range_id"),
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo",
			bytes.NewBuffer([]byte(`{"radioOn":false}`)))

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		r := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(*r.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(*r.ErrorMessage) > 0, "Response error message is missing")
		}
	}
}

func readResponse(w *httptest.ResponseRecorder) NetworksResponse {
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		panic(fmt.Sprintf("Error while reading body: %s", err))
	}

	var netsResponse NetworksResponse
	err = json.Unmarshal(bodyBytes, &netsResponse)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing body: %s", err))
	}

	return netsResponse
}
