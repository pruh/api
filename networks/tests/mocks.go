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
	MockGetWlans        func(omadaControllerId *string, loginToken *string,
		siteId *string) (*OmadaResponse, error)
	MockGetSsids func(omadaControllerId *string, loginToken *string,
		siteId *string, wlanId *string) (*OmadaResponse, error)
	MockUpdateSsid func(omadaControllerId *string, loginToken *string,
		siteId *string, wlanId *string, ssidId *string, scheduleId *string) (*OmadaResponse, error)
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

func (oa *MockOmadaApi) GetWlans(omadaControllerId *string, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return oa.MockGetWlans(omadaControllerId, loginToken, siteId)
}

func (oa *MockOmadaApi) GetSsids(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string) (*OmadaResponse, error) {
	return oa.MockGetSsids(omadaControllerId, loginToken, siteId, wlanId)
}

func (oa *MockOmadaApi) UpdateSsid(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string, ssidId *string, scheduleId *string) (*OmadaResponse, error) {
	return oa.MockUpdateSsid(omadaControllerId, loginToken, siteId, wlanId, ssidId, scheduleId)
}
