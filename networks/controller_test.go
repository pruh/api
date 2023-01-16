package networks_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	. "github.com/pruh/api/config/tests"
	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
)

func TestUpdateWifis_ControllerId(t *testing.T) {
	testsData := []struct {
		description        string
		requestUrl         string
		ssidParam          *string
		omadaResponseError bool
		omadaControllerId  *string
		loginToken         *string
		responseCode       int
	}{
		{
			description:       "happy path",
			requestUrl:        "https://omada.example.com/networks/ssid",
			ssidParam:         NewStr("my_ssid"),
			omadaControllerId: NewStr("c_id"),
			loginToken:        NewStr("login_token"),
			responseCode:      http.StatusOK,
		},
		{
			description:  "ssid missing in the request params",
			requestUrl:   "https://omada.example.com",
			responseCode: http.StatusBadRequest,
		},
		{
			description:        "omada controller id response error",
			requestUrl:         "https://omada.example.com",
			ssidParam:          NewStr("my_ssid"),
			omadaResponseError: true,
			omadaControllerId:  NewStr("c_id"),
			responseCode:       http.StatusBadGateway,
		},
		{
			description:       "controller id is missing in omada response",
			requestUrl:        "https://omada.example.com",
			ssidParam:         NewStr("my_ssid"),
			omadaControllerId: nil,
			responseCode:      http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(NewStr("8080"), NewStr("1"), NewStr("123"), nil, nil, nil, NewStr(testData.requestUrl), nil, nil),
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: testData.loginToken,
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
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
				MockGetWlans: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("ssid_id"),
									Name: NewStr("my_ssid"),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
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

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: testData.loginToken,
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("ssid_id"),
									Name: NewStr("my_ssid"),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
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
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("ssid_id"),
									Name: NewStr("my_ssid"),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{
								{
									Id:   NewStr("ssid_id"),
									Name: NewStr("my_ssid"),
								},
							},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					var res *Result
					if testData.includeSsids {
						res = &Result{
							Data: &[]Data{{Id: NewStr("ssid_id"), Name: NewStr("my_ssid")}},
						}
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result:    res,
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("ssid_id"), Name: NewStr("my_ssid")}},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {

					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
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
			description:  "happy path",
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
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Token: NewStr("login_token"),
						},
					}

					return resp, nil
				},
				MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
						},
					}

					return resp, nil
				},
				MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
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
				MockGetSsids: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
						Result: &Result{
							Data: &[]Data{{Id: NewStr("ssid_id"), Name: NewStr("my_ssid")}},
						},
					}

					return resp, nil
				},
				MockUpdateSsid: func(omadaControllerId *string, loginToken *string,
					siteId *string, wlanId *string, ssidId *string,
					scheduleId *string) (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}, nil
				},
				MockGetTimeRanges: func(omadaControllerId *string, loginToken *string,
					siteId *string) (*OmadaResponse, error) {

					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       NewStr("test"),
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "my_ssid",
		}
		req = mux.SetURLVars(req, vars)

		controller.UpdateWifi(w, req)

		netsResponse := readResponse(w)

		assert.Equal(testData.responseCode, w.Code, "Response code is not correct")
		if testData.responseCode == http.StatusOK {
			assert.True(netsResponse.Data.Updated, "Response success body missing updated flag")
		} else {
			assert.True(len(netsResponse.Error.Message) > 0, "Response error message is missing")
		}
	}
}

func readResponse(w *httptest.ResponseRecorder) NetworksResponse {
	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		panic(fmt.Sprintf("Error while reading body: %s", err))
	}

	var netsResponse NetworksResponse
	json.Unmarshal(bodyBytes, &netsResponse)
	return netsResponse
}
