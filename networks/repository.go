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
