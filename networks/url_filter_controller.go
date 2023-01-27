package networks

// Controller to operatate with URL filters
type urlFilterController struct {
	// config     *config.Configuration
	// ufc        UrlFilterController
	// repository Repository
}

// Controller to operatate with URL filters
type UrlFilterController interface {
	QueryUrlFilters(ssidData *Data) (*[]UrlFilter, error)
	MaybeUpdateUrlFilters() (*[]UrlFilter, error)
}

// creates new URL filter controller
func NewUrlFilterController() UrlFilterController {
	return urlFilterController{}
}

// Query URL filters for SSID
func (ufc urlFilterController) QueryUrlFilters(ssidData *Data) (*[]UrlFilter, error) {
	return nil, nil
}

// Update URL filter if required
func (ufc urlFilterController) MaybeUpdateUrlFilters() (*[]UrlFilter, error) {
	return nil, nil
}
