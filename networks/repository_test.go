package networks_test

import (
	"net/http"
	"testing"

	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
	"github.com/stretchr/testify/assert"
)

func TestRepoGetControllerId(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetControllerId: func() (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
				Result: &Result{
					OmadacId: NewStr("someId"),
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetControllerId()

	assert.True(mockCalled, "mock is not called")
	assert.Equal("someId", *controllerId.Result.OmadacId, "controller id is not as expected")
}

func TestRepoLogin(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockLogin: func(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
				Result: &Result{
					Token: NewStr("login_token"),
				},
			}

			cookies := []*http.Cookie{
				&http.Cookie{
					Name:  "cookie_name",
					Value: "cookie_value",
				},
			}

			mockCalled = true

			return resp, cookies, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, cookies, _ := repo.Login(nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal("cookie_name", cookies[0].Name, "login token is not as expected")
	assert.Equal("cookie_value", cookies[0].Value, "login token is not as expected")
	assert.Equal("login_token", *controllerId.Result.Token, "login token is not as expected")
}

func TestRepoGetSites(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetSites: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
				Result: &Result{
					Data: &[]Data{{Id: NewStr("site_id"), Name: NewStr("site_name")}},
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetSites(nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(Data{Id: NewStr("site_id"), Name: NewStr("site_name")},
		(*controllerId.Result.Data)[0], "sites are not as expected")
}

func TestRepoGetWlans(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetWlans: func(omadaControllerId *string, cookies []*http.Cookie,
			loginToken *string, siteId *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
				Result: &Result{
					Data: &[]Data{{Id: NewStr("wlan_id"), Name: NewStr("wlan_name")}},
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetWlans(nil, nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(Data{Id: NewStr("wlan_id"), Name: NewStr("wlan_name")},
		(*controllerId.Result.Data)[0], "wlans are not as expected")
}

func TestRepoGetSsids(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetSsids: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
			siteId *string, wlanId *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
				Result: &Result{
					Data: &[]Data{{Id: NewStr("ssid_id"), Name: NewStr("ssid_name")}},
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetSsids(nil, nil, nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(Data{Id: NewStr("ssid_id"), Name: NewStr("ssid_name")},
		(*controllerId.Result.Data)[0], "wlans are not as expected")
}

func TestRepoUpdateSsid(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockUpdateSsid: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
			siteId *string, wlanId *string, ssidId *string,
			ssidUpdateData *OmadaSsidUpdateData) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.UpdateSsid(nil, nil, nil, nil, nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(0, controllerId.ErrorCode, "wrong response data")
}

func TestRepoGetTimeRanges(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetTimeRanges: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
			siteId *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetTimeRanges(nil, nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(0, controllerId.ErrorCode, "wrong response data")
}

func TestRepoCreateTimeRanges(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockCreateTimeRange: func(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
			siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       NewStr("Success."),
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.CreateTimeRange(nil, nil, nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(0, controllerId.ErrorCode, "wrong response data")
}
