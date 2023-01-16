package networks_test

import (
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
				Msg:       StrPtr("Success."),
				Result: &Result{
					OmadacId: StrPtr("someId"),
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
		MockLogin: func(omadaControllerId *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       StrPtr("Success."),
				Result: &Result{
					Token: StrPtr("login_token"),
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.Login(nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal("login_token", *controllerId.Result.Token, "login token is not as expected")
}

func TestRepoGetSites(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetSites: func(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       StrPtr("Success."),
				Result: &Result{
					Data: &[]Data{{Id: StrPtr("site_id"), Name: StrPtr("site_name")}},
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetSites(nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(Data{Id: StrPtr("site_id"), Name: StrPtr("site_name")},
		(*controllerId.Result.Data)[0], "sites are not as expected")
}

func TestRepoGetWlans(t *testing.T) {
	var mockCalled = false

	mockOmadaApi := MockOmadaApi{
		MockGetWlans: func(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
			resp := &OmadaResponse{
				ErrorCode: 0,
				Msg:       StrPtr("Success."),
				Result: &Result{
					Data: &[]Data{{Id: StrPtr("wlan_id"), Name: StrPtr("wlan_name")}},
				},
			}

			mockCalled = true

			return resp, nil
		},
	}

	repo := NewRepository(&mockOmadaApi)

	assert := assert.New(t)

	controllerId, _ := repo.GetWlans(nil, nil, nil)

	assert.True(mockCalled, "mock is not called")
	assert.Equal(Data{Id: StrPtr("wlan_id"), Name: StrPtr("wlan_name")},
		(*controllerId.Result.Data)[0], "wlans are not as expected")
}
