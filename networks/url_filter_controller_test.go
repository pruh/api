package networks_test

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/pruh/api/networks"
	. "github.com/pruh/api/networks/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestControllerQueryUrlFilters(t *testing.T) {
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
					Name:   NewStr("Block List 1"),
					Enable: NewBool(true),
					Urls:   &[]string{"*google.com*", "*goo.gl*"},
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
		t.Logf("tesing %s", testData.description)

		ufc := NewUrlFilterController(NewRepository(&MockOmadaApi{
			MockQueryUrlFilters: func(omadaControllerId *string, cookies []*http.Cookie,
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
func TestMaybeUpdateUrlFilters(t *testing.T) {
	testsData := []struct {
		description          string
		ssidId               string
		requestFiltersUpdate *[]UrlFilter
		omadaResponse        *OmadaResponse
		expectUpdated        *bool
		expectFilters        *[]UrlFilter
		expectError          bool
	}{
		{
			description: "add single filter",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_2"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_2"},
							SourceIds:  &[]string{"my_ssid_id"},
						},
					},
				},
			},
			expectUpdated: NewBool(true),
			expectFilters: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
				{
					Name:   NewStr("filter_name_2"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_2"},
				},
			},
		},
		{
			description: "add when multiple ssid",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_1"},
							SourceIds:  &[]string{"another_ssid_1", "another_ssid_2"},
						},
					},
				},
			},
			expectUpdated: NewBool(true),
			expectFilters: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
		},
		{
			description: "update single filter",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_1"},
							SourceIds:  &[]string{"another_ssid_id"},
						},
					},
				},
			},
			expectUpdated: NewBool(true),
			expectFilters: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
		},
		{
			description: "delete single filter",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(false),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_1"},
							SourceIds:  &[]string{"my_ssid_id"},
						},
					},
				},
			},
			expectUpdated: NewBool(true),
			expectFilters: &[]UrlFilter{},
		},
		{
			description: "delete when several ssids",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(false),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_1"},
							SourceIds:  &[]string{"my_ssid_id", "another_ssid"},
						},
					},
				},
			},
			expectUpdated: NewBool(true),
			expectFilters: &[]UrlFilter{},
		},
		{
			description: "add not updated",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{
						{
							Id:         NewStr("1"),
							Name:       NewStr("filter_name_1"),
							SiteId:     NewStr("site_id"),
							Policy:     NewInt(ENABLE_FILTERING),
							EntryId:    NewInt(123),
							Status:     NewBool(true),
							SourceType: NewInt(SSID_SOURCE_TYPE),
							Urls:       &[]string{"url_1"},
							SourceIds:  &[]string{"my_ssid_id"},
						},
					},
				},
			},
			expectUpdated: NewBool(false),
			expectFilters: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(true),
					Urls:   &[]string{"url_1"},
				},
			},
		},
		{
			description: "delete not updated",
			ssidId:      "my_ssid_id",
			requestFiltersUpdate: &[]UrlFilter{
				{
					Name:   NewStr("filter_name_1"),
					Enable: NewBool(false),
					Urls:   &[]string{"url_1"},
				},
			},
			omadaResponse: &OmadaResponse{
				Result: &Result{
					Data: &[]Data{},
				},
			},
			expectUpdated: NewBool(false),
			expectFilters: &[]UrlFilter{},
		},
		{
			description: "error if upstream HTTP error",
			expectError: true,
		},
	}

	assert := assert.New(t)

	for _, testData := range testsData {
		t.Logf("tesing %s", testData.description)

		ufc := NewUrlFilterController(NewRepository(&MockOmadaApi{
			MockQueryUrlFilters: func(omadaControllerId *string, cookies []*http.Cookie,
				loginToken, siteId *string) (*OmadaResponse, error) {
				if testData.expectError {
					return nil, errors.New("test")
				}

				return testData.omadaResponse, nil
			},
			MockCreateUrlFilter: func(omadaControllerId *string, cookies []*http.Cookie, loginToken,
				siteId *string, urlFilterData *Data) (*OmadaResponse, error) {
				if testData.expectError {
					return nil, errors.New("test")
				}

				return &OmadaResponse{}, nil
			},
			MockUpdateUrlFilter: func(omadaControllerId *string, cookies []*http.Cookie, loginToken,
				siteId *string, urlFilterData *Data) (*OmadaResponse, error) {
				if testData.expectError {
					return nil, errors.New("test")
				}

				return &OmadaResponse{}, nil
			},
			MockDeleteUrlFilter: func(omadaControllerId *string, cookies []*http.Cookie,
				loginToken, siteId, urlFilterId *string) (*OmadaResponse, error) {
				if testData.expectError {
					return nil, errors.New("test")
				}

				return &OmadaResponse{}, nil
			},
		}))

		urlFilters, updated, err := ufc.MaybeUpdateUrlFilters(NewStr("cid"), []*http.Cookie{},
			NewStr("login_token"), NewStr("site_id"), &Data{Id: &testData.ssidId}, testData.requestFiltersUpdate)
		if testData.expectError {
			assert.True(err != nil, "should return error")
			continue
		}

		assert.Equal(*testData.expectUpdated, *updated, "updated flag is incorrect")
		compareFunc := func(a, b UrlFilter) bool {
			return *a.Name < *b.Name
		}
		slices.SortFunc(*testData.expectFilters, compareFunc)
		slices.SortFunc(*urlFilters, compareFunc)
		assert.Equal(*testData.expectFilters, *urlFilters, "returned filters are incorrect")
	}
}
