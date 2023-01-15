package networks_test

import (
	"errors"
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

		assert.Equal(testData.omadaErrorCode, omadaControllerId.ErrorCode, "Response code is not correct")
		assert.Equal(testData.omadacId, omadaControllerId.Result.OmadacId, "omada id parsing is not correct")
	}
}
