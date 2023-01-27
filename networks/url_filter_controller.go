package networks

// Controller to operatate with URL filters
type urlFilterController struct {
	repository Repository
}

// Controller to operatate with URL filters
type UrlFilterController interface {
	QueryUrlFilters(ssidData *Data) (*[]UrlFilter, error)
	MaybeUpdateUrlFilters() (*[]UrlFilter, error)
}

// creates new URL filter controller
func NewUrlFilterController(r Repository) UrlFilterController {
	return urlFilterController{
		repository: r,
	}
}

// Query URL filters for SSID
func (ufc urlFilterController) QueryUrlFilters(ssidData *Data) (*[]UrlFilter, error) {
	// todo query omada for all AP url filters

	// extract uf for given ssid
	return nil, nil
}

// Update URL filter if required
func (ufc urlFilterController) MaybeUpdateUrlFilters() (*[]UrlFilter, error) {
	return nil, nil
}
