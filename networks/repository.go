package networks

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

func (r *Repository) Login(omadaControllerId *string) (*OmadaResponse, error) {
	return r.omadaApi.Login(omadaControllerId)
}

func (r *Repository) GetSites(omadaControllerId *string,
	loginToken *string) (*OmadaResponse, error) {
	return r.omadaApi.GetSites(omadaControllerId, loginToken)
}

func (r *Repository) GetWlans(omadaControllerId *string, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	return r.omadaApi.GetWlans(omadaControllerId, loginToken, siteId)
}

func (r *Repository) GetSsids(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string) (*OmadaResponse, error) {
	return r.omadaApi.GetSsids(omadaControllerId, loginToken, siteId, wlanId)
}

func (r *Repository) UpdateSsid(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string, ssidId *string, scheduleId *string) (*OmadaResponse, error) {
	return r.omadaApi.UpdateSsid(omadaControllerId, loginToken, siteId, wlanId, ssidId, scheduleId)
}
