package networks

import (
	"net/http"

	"github.com/golang/glog"
	"golang.org/x/exp/maps"
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
		urlFilter.Enable = NewBool(*omadaFilter.Policy == ENABLE_FILTERING)
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
		glog.Info("no need to set any filters, just return whatever we have")
		filters, err := ufc.QueryUrlFilters(omadaControllerId, cookies, loginToken, siteId, ssidData)
		if err != nil {
			return nil, nil, err
		}

		return filters, NewBool(false), nil
	}

	glog.Info("query existing filters")
	resp, err := ufc.repository.QueryUrlFilters(omadaControllerId, cookies, loginToken, siteId)
	if err != nil {
		return nil, nil, err
	}

	updated := false

	glog.Info("find filters that have the same urls")
	createdFilters := []Data{}
	deletedFilters := []Data{}
	for _, requestedFilter := range *requestedFilters {
		reqFilterUrls := *requestedFilter.Urls
		slices.Sort(reqFilterUrls)

		if *requestedFilter.Enable {
			f, u, err := ufc.addFiltering(omadaControllerId, cookies, loginToken, siteId,
				requestedFilter, reqFilterUrls, ssidData, resp.Result.Data)
			if err != nil {
				return nil, nil, err
			}
			updated = updated || *u

			createdFilters = append(createdFilters, *f...)
		} else {
			f, u, err := ufc.deleteFiltering(omadaControllerId, cookies, loginToken, siteId,
				requestedFilter, reqFilterUrls, ssidData, resp.Result.Data)
			if err != nil {
				return nil, nil, err
			}
			updated = updated || *u

			deletedFilters = append(deletedFilters, *f...)
		}
	}

	newFilters := map[string]Data{}
	for _, f := range *resp.Result.Data {
		newFilters[*f.Id] = f
	}
	// delete deleted
	for _, f := range deletedFilters {
		delete(newFilters, *f.Id)
	}
	nf := maps.Values(newFilters)
	// add newly created filters
	nf = append(nf, createdFilters...)

	glog.Infof("new filters %+v", nf)

	return convertOmadaUrlFilters(&nf, ssidData.Id), NewBool(updated), nil
}

func (ufc urlFilterController) deleteFiltering(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string,
	requestedFilter UrlFilter, reqFilterUrls []string,
	ssidData *Data, omadaFilters *[]Data) (*[]Data, *bool, error) {

	glog.Infof("looking to delete filter %+v", requestedFilter)

	deletedFilters := []Data{}
	updated := false
	for _, omadaFilter := range *omadaFilters {
		if *omadaFilter.SourceType != SSID_SOURCE_TYPE {
			glog.Info("not ssid filter")
			continue
		}

		if !*omadaFilter.Status || *omadaFilter.Policy != ENABLE_FILTERING ||
			*omadaFilter.Name != *requestedFilter.Name {
			glog.Info("filter does not match")
			continue
		}

		omadaFilterUrls := *omadaFilter.Urls
		slices.Sort(omadaFilterUrls)

		if !slices.Equal(reqFilterUrls, omadaFilterUrls) {
			glog.Info("urls do not match")
			continue
		}

		if !slices.Contains(*omadaFilter.SourceIds, *ssidData.Id) {
			glog.Info("filtering is not for requested ssid")
			continue
		}

		if len(*omadaFilter.SourceIds) > 1 {
			glog.Info("multiple ssids in filtering, deleleting current ssid")
			i := slices.Index(*omadaFilter.SourceIds, *ssidData.Id)
			slices.Delete(*omadaFilter.SourceIds, i, i+1)
			_, err := ufc.repository.UpdateUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, &omadaFilter)
			if err != nil {
				glog.Info("upstream error updating a rule")
				return nil, nil, err
			}

			updated = true
		} else {
			glog.Info("single ssids in filtering, deleleting the rule")
			_, err := ufc.repository.DeleteUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, omadaFilter.Id)
			if err != nil {
				glog.Info("upstream error deleting a rule")
				return nil, nil, err
			}

			updated = true

			deletedFilters = append(deletedFilters, omadaFilter)
		}
	}

	return &deletedFilters, NewBool(updated), nil
}

func (ufc urlFilterController) addFiltering(omadaControllerId *string,
	cookies []*http.Cookie, loginToken *string, siteId *string,
	requestedFilter UrlFilter, reqFilterUrls []string,
	ssidData *Data, omadaFilters *[]Data) (*[]Data, *bool, error) {

	glog.Infof("looking to add filter %+v", requestedFilter)

	createdFilters := []Data{}
	updated := false
	blockingRuleFound := false
	for _, omadaFilter := range *omadaFilters {
		if *omadaFilter.SourceType != SSID_SOURCE_TYPE {
			glog.Info("not ssid filter")
			continue
		}

		if !*omadaFilter.Status || *omadaFilter.Policy != ENABLE_FILTERING ||
			*omadaFilter.Name != *requestedFilter.Name {
			glog.Info("filter does not match")
			continue
		}

		omadaFilterUrls := *omadaFilter.Urls
		slices.Sort(omadaFilterUrls)

		if !slices.Equal(reqFilterUrls, omadaFilterUrls) {
			glog.Info("urls do not match")
			continue
		}

		blockingRuleFound = true
		if !slices.Contains(*omadaFilter.SourceIds, *ssidData.Id) {
			glog.Info("no ssid in filtering, adding it")

			*omadaFilter.SourceIds = append(*omadaFilter.SourceIds, *ssidData.Id)

			_, err := ufc.repository.UpdateUrlFilter(omadaControllerId, cookies,
				loginToken, siteId, &omadaFilter)

			if err != nil {
				glog.Info("upstream error updating a rule")
				return nil, nil, err
			}

			updated = true
		}
	}

	if !blockingRuleFound {
		glog.Info("no blocking rule found, creating a new one")
		newFilter := Data{
			Name:       requestedFilter.Name,
			Status:     NewBool(true),
			Policy:     NewInt(ENABLE_FILTERING),
			SourceType: NewInt(SSID_SOURCE_TYPE),
			SourceIds:  &[]string{*ssidData.Id},
			Urls:       requestedFilter.Urls,
		}
		_, err := ufc.repository.CreateUrlFilter(omadaControllerId, cookies,
			loginToken, siteId, &newFilter)

		if err != nil {
			glog.Info("upstream error creating a new rule")
			return nil, nil, err
		}

		updated = true

		createdFilters = append(createdFilters, newFilter)
	}

	return &createdFilters, NewBool(updated), nil
}
