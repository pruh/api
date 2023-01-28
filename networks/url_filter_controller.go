package networks

import (
	"net/http"

	"golang.org/x/exp/slices"
)

const ENABLE_FILTERING = 0
const BYPASS_FILTERING = 1
const SSID_SOURCE_TYPE = 2

// Controller to operatate with URL filters
type urlFilterController struct {
	repository Repository
}

// Controller to operatate with URL filters
type UrlFilterController interface {
	QueryUrlFilters(omadaControllerId *string, cookies []*http.Cookie,
		loginToken *string, siteId *string, ssidData *Data) (*[]UrlFilter, error)
	MaybeUpdateUrlFilters(omadaControllerId *string,
		cookies []*http.Cookie, loginToken *string, siteId *string, ssidData *Data,
		requestedFilters *[]UrlFilter) (*[]UrlFilter, *bool, error)
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
	resp, err := ufc.repository.QueryUrlFilters(omadaControllerId, cookies, loginToken, siteId)
	if err != nil {
		return nil, err
	}

	return convertOmadaUrlFilters(resp.Result.Data, ssidData.Id), nil
}

func convertOmadaUrlFilters(omadaFilters *[]Data, ssidId *string) *[]UrlFilter {
	urlFilters := []UrlFilter{}
	for _, omadaFilter := range *omadaFilters {
		if *omadaFilter.SourceType != SSID_SOURCE_TYPE {
			// only accepct SSID filters
			continue
		}
		if *omadaFilter.SourceIds == nil ||
			!slices.Contains(*omadaFilter.SourceIds, *ssidId) {
			// skip filters that do not belong to requested ssid
			continue
		}
		var urlFilter UrlFilter
		urlFilter.Name = omadaFilter.Name
		urlFilter.BypassFilter = NewBool(*omadaFilter.Policy == BYPASS_FILTERING)
		urlFilter.Urls = omadaFilter.Urls

		urlFilters = append(urlFilters, urlFilter)
	}

	return &urlFilters
}

// Update URL filter if required
func (ufc urlFilterController) MaybeUpdateUrlFilters(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string, ssidData *Data,
	requestedFilters *[]UrlFilter) (*[]UrlFilter, *bool, error) {

	if requestedFilters == nil {
		// no need to set any filters, just return whatever we have
		filters, err := ufc.QueryUrlFilters(omadaControllerId, cookies, loginToken, siteId, ssidData)
		if err != nil {
			return nil, nil, err
		}

		return filters, NewBool(false), nil
	}

	// query existing filters
	resp, err := ufc.repository.QueryUrlFilters(omadaControllerId, cookies, loginToken, siteId)
	if err != nil {
		return nil, nil, err
	}

	updated := false

	// find filters that have the same urls
	for _, requestedFilter := range *requestedFilters {
		reqFilterUrls := *requestedFilter.Urls
		slices.Sort(reqFilterUrls)

		if *requestedFilter.BypassFilter {
			ufc.removeFiltering(omadaControllerId, cookies, loginToken, siteId,
				requestedFilter, reqFilterUrls, ssidData, resp.Result.Data)
		} else {
			ufc.addFiltering(omadaControllerId, cookies, loginToken, siteId,
				requestedFilter, reqFilterUrls, ssidData, resp.Result.Data)
		}
	}

	return convertOmadaUrlFilters(resp.Result.Data, ssidData.Id), NewBool(updated), nil
}

func (ufc urlFilterController) removeFiltering(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string,
	requestedFilter UrlFilter, reqFilterUrls []string,
	ssidData *Data, omadaFilters *[]Data) (*bool, error) {

	updated := false
	for _, omadaFilter := range *omadaFilters {
		if *omadaFilter.SourceType != SSID_SOURCE_TYPE {
			// only accepct SSID filters
			continue
		}

		if *omadaFilter.Status == false || *omadaFilter.Policy != ENABLE_FILTERING {
			// rule is not enabled or policy is not blocking
			continue
		}

		omadaFilterUrls := *omadaFilter.Urls
		slices.Sort(omadaFilterUrls)

		if !slices.Equal(reqFilterUrls, omadaFilterUrls) {
			// urls do not match
			continue
		}

		if !slices.Contains(*omadaFilter.SourceIds, *ssidData.Id) {
			// rule is not for requested ssid
			continue
		}

		if len(*omadaFilter.SourceIds) > 1 {
			i := slices.Index(*omadaFilter.SourceIds, *ssidData.Id)
			slices.Delete(*omadaFilter.SourceIds, i, i+1)
			_, err := ufc.repository.UpdateUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, &omadaFilter)
			if err != nil {
				return nil, err
			}

			updated = true
		} else {
			_, err := ufc.repository.DeleteUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, omadaFilter.Id)
			if err != nil {
				return nil, err
			}

			updated = true
		}
	}

	return NewBool(updated), nil
}

func (ufc urlFilterController) addFiltering(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string,
	requestedFilter UrlFilter, reqFilterUrls []string,
	ssidData *Data, omadaFilters *[]Data) (*bool, error) {

	updated := false
	blockingRuleFound := false
	for _, omadaFilter := range *omadaFilters {
		if *omadaFilter.SourceType != SSID_SOURCE_TYPE {
			// only accepct SSID filters
			continue
		}

		if *omadaFilter.Status == false || *omadaFilter.Policy != ENABLE_FILTERING {
			// rule is not enabled or policy is not blocking
			continue
		}

		omadaFilterUrls := *omadaFilter.Urls
		slices.Sort(omadaFilterUrls)

		if !slices.Equal(reqFilterUrls, omadaFilterUrls) {
			// urls do not match
			continue
		}

		blockingRuleFound = true
		if !slices.Contains(*omadaFilter.SourceIds, *ssidData.Id) {
			*omadaFilter.SourceIds = append(*omadaFilter.SourceIds, *ssidData.Id)

			_, err := ufc.repository.UpdateUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, &omadaFilter)

			if err != nil {
				return nil, err
			}

			updated = true
		}
	}

	if !blockingRuleFound {
		_, err := ufc.repository.CreateUrlFilter(omadaControllerId, cookies,
			loginToken, siteId, &Data{
				Name:       requestedFilter.Name,
				Status:     NewBool(true),
				Policy:     NewInt(ENABLE_FILTERING),
				SourceType: NewInt(SSID_SOURCE_TYPE),
				SourceIds:  &[]string{*ssidData.Id},
				Urls:       requestedFilter.Urls,
			},
		)

		if err != nil {
			return nil, err
		}

		updated = true
	}

	return NewBool(updated), nil
}
