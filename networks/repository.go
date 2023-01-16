package networks

import "net/http"

// Repository to interact with networks data providers
type Repository struct {
	omadaApi OmadaApi
}

// NewRepository creates new networks repository
func NewRepository(omadaApi OmadaApi) Repository {
	return Repository{
		omadaApi: omadaApi,
	}
}

func (r *Repository) GetControllerId() (*OmadaResponse, error) {
	return r.omadaApi.GetControllerId()
}

func (r *Repository) Login(omadaControllerId *string) (*OmadaResponse, []*http.Cookie, error) {
	return r.omadaApi.Login(omadaControllerId)
}

func (r *Repository) GetSites(omadaControllerId *string, cookies []*http.Cookie,
	loginToken *string) (*OmadaResponse, error) {
	return r.omadaApi.GetSites(omadaControllerId, cookies, loginToken)
}

func (r *Repository) GetWlans(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return r.omadaApi.GetWlans(omadaControllerId, cookies, loginToken, siteId)
}

func (r *Repository) GetSsids(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, wlanId *string) (*OmadaResponse, error) {
	return r.omadaApi.GetSsids(omadaControllerId, cookies, loginToken, siteId, wlanId)
}

func (r *Repository) UpdateSsid(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, wlanId *string, ssidId *string,
	ssidUpdateData *OmadaSsidUpdateData) (*OmadaResponse, error) {
	return r.omadaApi.UpdateSsid(omadaControllerId, cookies, loginToken, siteId, wlanId, ssidId,
		ssidUpdateData)
}

func (r *Repository) GetTimeRanges(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return r.omadaApi.GetTimeRanges(omadaControllerId, cookies, loginToken, siteId)
}

func (r *Repository) CreateTimeRange(omadaControllerId *string, cookies []*http.Cookie, loginToken *string,
	siteId *string, timeRangeData *Data) (*OmadaResponse, error) {
	return r.omadaApi.CreateTimeRange(omadaControllerId, cookies, loginToken, siteId, timeRangeData)
}