package networks_test

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
	"github.com/stretchr/testify/assert"
)

func TestQueryUrlFilters(t *testing.T) {
	testsData := []struct {
		description      string
		ssidId           string
		omadaResponse    *OmadaResponse
		expectUrlFilters *[]UrlFilter
		expectError      bool
	}{
		{
			description: "one filter",
			ssidId:      "my_ssid_id",
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("Block List 1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(0),
							Type:       NewStr("ap"),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(2),
							Urls:       &[]string{"*google.com*", "*goo.gl*"},
							SourceIds:  &[]string{"my_ssid_id", "another_ssid_id"},
						},
					},
				},
			},
			expectUrlFilters: &[]UrlFilter{
				{
					Name:         NewStr("Block List 1"),
					BypassFilter: NewBool(false),
					Urls:         &[]string{"*google.com*", "*goo.gl*"},
				},
			},
		},
		{
			description: "three filters, wrong ssid",
			ssidId:      "my_ssid_id",
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("Block List 1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(0),
							Type:       NewStr("ap"),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(2),
							Urls:       &[]string{"*google.com*", "*goo.gl*"},
							SourceIds:  &[]string{"another_ssid_id"},
						},
						{
							Id:         NewStr("1"),
							Name:       NewStr("Block List 2"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(0),
							Type:       NewStr("ap"),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(2),
							Urls:       &[]string{"*google.com*", "*goo.gl*"},
							SourceIds:  &[]string{"another_ssid_id"},
						},
						{
							Id:         NewStr("1"),
							Name:       NewStr("Block List 3"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(0),
							Type:       NewStr("ap"),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(2),
							Urls:       &[]string{"*google.com*", "*goo.gl*"},
							SourceIds:  &[]string{"another_ssid_id"},
						},
					},
				},
			},
			expectUrlFilters: &[]UrlFilter{},
		},
		{
			description: "error if upstream HTTP error",
			expectError: true,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %+v", testData.description)

		ufc := NewUrlFilterController(NewRepository(&MockOmadaApi{
			MockQueryAPUrlFilters: func(omadaControllerId *string, cookies []*http.Cookie,
				loginToken, siteId *string) (*OmadaResponse, error) {
				if testData.expectError {
					return nil, errors.New("test")
				}

				return testData.omadaResponse, nil
			},
		}))

		urlFilters, err := ufc.QueryUrlFilters(NewStr("cid"), []*http.Cookie{},
			NewStr("login_token"), NewStr("site_id"), &Data{Id: &testData.ssidId})
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(*testData.expectUrlFilters, *urlFilters, "Error code is not correct")
	}
}
