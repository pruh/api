package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/pruh/api/http"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	testsData := []struct {
		description  string
		responseCode int
		responseBody string
	}{
		{
			description:  "normal config",
			responseCode: 200,
			responseBody: "OK",
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(testData.responseBody))
			if err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		client := NewHTTPClient()
		r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(r)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(testData.responseCode, resp.StatusCode)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(testData.responseBody, string(respBody))
	}
}
