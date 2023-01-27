package networks

import "net/http"

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
	// todo query omada for all AP url filters
	ufc.repository.QueryAPUrlFilters(omadaControllerId, cookies, loginToken, siteId)

	// extract uf for given ssid
	return nil, nil
}

// Update URL filter if required
func (ufc urlFilterController) MaybeUpdateUrlFilters() (*[]UrlFilter, error) {
	return nil, nil
}
