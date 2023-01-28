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
	MockLogin           func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error)
	MockGetSites        func(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string) (*OmadaResponse, error)
	MockGetWlans func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
		siteId *string) (*OmadaResponse, error)
	MockGetSsids func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
		siteId *string, wlanId *string) (*OmadaResponse, error)
	MockUpdateSsid func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
		siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error)
	MockGetTimeRanges func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
		siteId *string) (*OmadaResponse, error)
	MockCreateTimeRange func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
		siteId *string, trData *Data) (*OmadaResponse, error)

	MockQueryUrlFilters func(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string) (*OmadaResponse, error)
	MockCreateUrlFilter func(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string, urlFilterData *Data) (*OmadaResponse, error)
	MockDeleteUrlFilter func(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string, urlFilterId *string) (*OmadaResponse, error)
}

func (oa *MockOmadaApi) GetControllerId() (*OmadaResponse, error) {
	return oa.MockGetControllerId()
}

func (oa *MockOmadaApi) Login(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
	return oa.MockLogin(omadaControllerId)
}

func (oa *MockOmadaApi) GetSites(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string) (*OmadaResponse, error) {
	return oa.MockGetSites(omadaControllerId, cookies, loginToken)
}

func (oa *MockOmadaApi) GetWlans(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return oa.MockGetWlans(omadaControllerId, cookies, loginToken, siteId)
}

func (oa *MockOmadaApi) GetSsids(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, wlanId *string) (*OmadaResponse, error) {
	return oa.MockGetSsids(omadaControllerId, cookies, loginToken, siteId, wlanId)
}

func (oa *MockOmadaApi) UpdateSsid(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, wlanId *string, ssidUpdateData *Data) (*OmadaResponse, error) {
	return oa.MockUpdateSsid(omadaControllerId, cookies, loginToken, siteId, wlanId, ssidUpdateData)
}

func (oa *MockOmadaApi) GetTimeRanges(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return oa.MockGetTimeRanges(omadaControllerId, cookies, loginToken, siteId)
}

func (oa *MockOmadaApi) CreateTimeRange(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, trData *Data) (*OmadaResponse, error) {
	return oa.MockCreateTimeRange(omadaControllerId, cookies, loginToken, siteId, trData)
}

func (oa *MockOmadaApi) QueryUrlFilters(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return oa.MockQueryUrlFilters(omadaControllerId, cookies, loginToken, siteId)
}

func (oa *MockOmadaApi) CreateUrlFilter(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string, siteId *string, urlFilterData *Data) (*OmadaResponse, error) {
	return oa.MockCreateUrlFilter(omadaControllerId, cookies, loginToken, siteId, urlFilterData)
}

func (oa *MockOmadaApi) DeleteUrlFilter(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string, siteId *string, urlFilterId *string) (*OmadaResponse, error) {
	return oa.MockDeleteUrlFilter(omadaControllerId, cookies, loginToken, siteId, urlFilterId)
}

type MockUrlFilterController struct {
	MockQueryUrlFilters func(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error)
	MockMaybeUpdateUrlFilters func(omadaControllerId *string,
		cookies []*http.Cookie, loginToken *string, siteId *string, ssidData *Data,
		requestedFilters *[]UrlFilter) (*[]UrlFilter, *bool, error)
}

func (ufc MockUrlFilterController) QueryUrlFilters(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error) {
	return ufc.MockQueryUrlFilters(omadaControllerId, cookies, loginToken, siteId, ssidData)
}

func (ufc MockUrlFilterController) MaybeUpdateUrlFilters(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string, ssidData *Data,
	requestedFilters *[]UrlFilter) (*[]UrlFilter, *bool, error) {
	return ufc.MockMaybeUpdateUrlFilters(omadaControllerId,
		cookies, loginToken, siteId, ssidData, requestedFilters)
}

func NewMockUrlFilterController() MockUrlFilterController {
	return MockUrlFilterController{
		MockQueryUrlFilters: func(omadaControllerId *string, cookies []*http.Cookie,
			loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error) {
			return &[]UrlFilter{
				{
					Name:         NewStr("test"),
					BypassFilter: NewBool(false),
					Urls:         &[]string{"test_url"},
				},
			}, nil
		},
		MockMaybeUpdateUrlFilters: func(omadaControllerId *string, cookies []*http.Cookie,
			loginToken *string, siteId *string, ssidData *Data,
			requestedFilters *[]UrlFilter) (*[]UrlFilter, *bool, error) {
			return &[]UrlFilter{
				{
					Name:         NewStr("test"),
					BypassFilter: NewBool(false),
					Urls:         &[]string{"test_url"},
				},
			}, NewBool(false), nil
		},
	}
}
