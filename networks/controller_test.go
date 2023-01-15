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
			ssidParam:         StrPtr("my_ssid"),
			omadaControllerId: StrPtr("c_id"),
			loginToken:        StrPtr("login_token"),
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
			ssidParam:          StrPtr("my_ssid"),
			omadaResponseError: true,
			omadaControllerId:  StrPtr("c_id"),
			responseCode:       http.StatusBadGateway,
		},
		{
			description:       "controller id is missing in omada response",
			requestUrl:        "https://omada.example.com",
			ssidParam:         StrPtr("my_ssid"),
			omadaControllerId: nil,
			responseCode:      http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil, StrPtr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       "test",
						Result: Result{
							OmadacId: testData.omadaControllerId,
						},
					}

					return resp, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       "test",
						Result: Result{
							Token: testData.loginToken,
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
			loginToken:   StrPtr("login_token"),
			responseCode: http.StatusOK,
		},
		{
			description:        "omada controller id response error",
			requestUrl:         "https://omada.example.com",
			loginToken:         StrPtr("login_token"),
			omadaResponseError: true,
			responseCode:       http.StatusBadGateway,
		},
		{
			description:  "controller id is missing in omada response",
			requestUrl:   "https://omada.example.com",
			loginToken:   nil,
			responseCode: http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil, StrPtr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*OmadaResponse, error) {
					return &OmadaResponse{
						ErrorCode: 0,
						Msg:       "test",
						Result: Result{
							OmadacId: StrPtr("omada_cid"),
						},
					}, nil
				},
				MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					resp := &OmadaResponse{
						ErrorCode: 0,
						Msg:       "test",
						Result: Result{
							Token: testData.loginToken,
						},
					}

					return resp, nil
				},
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://example.com/foo", nil)

		// setting mux vars for testing
		vars := map[string]string{
			"ssid": "some_ssid",
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
