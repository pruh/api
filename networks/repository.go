package networks

// Repository to interact with networks data providers
type Repository struct {
	OmadaApi OmadaApi
}

// NewRepository creates new networks repository
func NewRepository(omadaApi OmadaApi) Repository {
	return Repository{
		OmadaApi: omadaApi,
	}
}

func (r *Repository) GetControllerId() (*ControllerIdResponse, error) {
	return r.OmadaApi.GetControllerId()
}
