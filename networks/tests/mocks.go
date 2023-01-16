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
	MockGetControllerId func() (*OmadaResponse, error)
	MockLogin           func(omadaControllerId *string) (*OmadaResponse, error)
	MockGetSites        func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error)
	MockGetWlans        func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error)
}

func (oa *MockOmadaApi) GetControllerId() (*OmadaResponse, error) {
	return oa.MockGetControllerId()
}

func (oa *MockOmadaApi) Login(omadaControllerId *string) (*OmadaResponse, error) {
	return oa.MockLogin(omadaControllerId)
}

func (oa *MockOmadaApi) GetSites(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
	return oa.MockGetSites(omadaControllerId, loginToken)
}

func (oa *MockOmadaApi) GetWlans(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
	return oa.MockGetWlans(omadaControllerId, loginToken, siteId)
}

func StrPtr(str string) *string {
	return &str
}

func IntPtr(num int) *int {
	return &num
}
