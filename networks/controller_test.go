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

func TestUpdateWifis(t *testing.T) {
	testsData := []struct {
		description        string
		requestUrl         string
		ssidParam          *string
		omadaResponseError bool
		omadaControllerId  *string
		responseCode       int
	}{
		{
			description:  "happy path",
			requestUrl:   "https://omada.example.com/networks/ssid",
			ssidParam:    StrPtr("my_ssid"),
			responseCode: http.StatusOK,
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
			responseCode:       http.StatusBadGateway,
		},
		{
			description:       "controller id is missing in omada response",
			requestUrl:        "https://omada.example.com",
			ssidParam:         StrPtr("my_ssid"),
			omadaControllerId: StrPtr(""),
			responseCode:      http.StatusBadGateway,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		controller := NewControllerWithParams(
			NewConfigSafe(StrPtr("8080"), StrPtr("1"), StrPtr("123"), nil, nil, nil, StrPtr(testData.requestUrl), nil, nil),
			&MockOmadaApi{
				MockGetControllerId: func() (*ControllerIdResponse, error) {
					if testData.omadaResponseError {
						return nil, errors.New("test")
					}

					var cId string
					if testData.omadaControllerId != nil {
						cId = *testData.omadaControllerId
					} else {
						cId = "some_id"
					}

					resp := &ControllerIdResponse{
						ErrorCode: 0,
						Msg:       "test",
						Result: ControllerIdResult{
							OmadacId: cId,
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

		controller.UpdateWifis(w, req)

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
