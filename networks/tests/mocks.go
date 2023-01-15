package tests

import (
	"net/http"

	. "github.com/pruh/api/networks"
)

type MockHTTPClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.MockDo(req)
}

type MockOmadaApi struct {
	MockGetControllerId func() (*ControllerIdResponse, error)
}

func (oa *MockOmadaApi) GetControllerId() (*ControllerIdResponse, error) {
	return oa.MockGetControllerId()
}

func StrPtr(str string) *string {
	return &str
}

func IntPtr(num int) *int {
	return &num
}
