package networks

import (
	"net/http"

	"golang.org/x/exp/slices"
)

// Controller to operatate with URL filters
type urlFilterController struct {
	repository Repository
}

// Controller to operatate with URL filters
type UrlFilterController interface {
	QueryUrlFilters(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error)
	MaybeUpdateUrlFilters() (*[]UrlFilter, error)
}

// creates new URL filter controller
func NewUrlFilterController(r Repository) UrlFilterController {
	return urlFilterController{
		repository: r,
	}
}

// Query URL filters for SSID
func (ufc urlFilterController) QueryUrlFilters(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error) {
	resp, err := ufc.repository.QueryAPUrlFilters(omadaControllerId, cookies, loginToken, siteId)
	if err != nil {
		return nil, err
	}

	urlFilters := []UrlFilter{}
	for _, omadaFilter := range *resp.Result.Data {
		if *omadaFilter.SourceType != 2 {
			// only accepct SSID filters
			continue
		}
		if *omadaFilter.SourceIds == nil ||
			!slices.Contains(*omadaFilter.SourceIds, *ssidData.Id) {
			// skip filters that do not belong to requested ssid
			continue
		}
		var urlFilter UrlFilter
		urlFilter.Name = omadaFilter.Name
		urlFilter.BypassFilter = NewBool(*omadaFilter.Policy == 1)
		urlFilter.Urls = omadaFilter.Urls

		urlFilters = append(urlFilters, urlFilter)
	}

	return &urlFilters, nil
}

// Update URL filter if required
func (ufc urlFilterController) MaybeUpdateUrlFilters() (*[]UrlFilter, error) {
	return nil, nil
}
