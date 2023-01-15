package networks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pruh/api/config"
	apihttp "github.com/pruh/api/http"
)

const (
	// networks path prefix
	networksPrefix = "/networks"

	// Wifis return list of sites for network controller
	Wifis = networksPrefix + "/wifis/{ssid}"
)

// Controller handles all network related requests
type controller struct {
	config     *config.Configuration
	repository Repository
}

type Controller interface {
	UpdateWifi(w http.ResponseWriter, r *http.Request)
}

// NewController creates new networks controller
func NewController(config *config.Configuration) Controller {
	omadaApi := NewOmadaApi(config, apihttp.NewHTTPClient())
	return NewControllerWithParams(config, omadaApi)
}

// Creates a new controller with additional dependencies for tests
func NewControllerWithParams(config *config.Configuration, omadaApi OmadaApi) Controller {
	return &controller{
		config:     config,
		repository: NewRepository(omadaApi),
	}
}

func (c *controller) UpdateWifi(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		errorMessage := "ssid is missing in the request parameters"
		c.writeResponse(w, http.StatusBadRequest, false, &errorMessage)
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		errorMessage := fmt.Sprintf("Omada Controller Id Query Error: %+v", err)
		c.writeResponse(w, http.StatusBadGateway, false, &errorMessage)
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaIdResp.ErrorCode != 0 || omadaLoginResp.Result.Token == nil {
		errorMessage := fmt.Sprintf("Omada Controller Id Query Error: %+v", err)
		c.writeResponse(w, http.StatusBadGateway, false, &errorMessage)
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId,
		omadaLoginResp.Result.Token)
	if err != nil || omadaIdResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		errorMessage := fmt.Sprintf("Omada Sites Query Error: %+v", err)
		c.writeResponse(w, http.StatusBadGateway, false, &errorMessage)
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	// get site id
	// and use the first one
	// GET /{omadacId}/api/v2/sites

	// get wlan id
	// and use the frist one
	// GET /{omadacId}/api/v2/sites/{siteId}/setting/wlans

	// get list of all ssids
	// find ssid that matches passed ssid
	// GET /{omadacId}/api/v2/sites/{siteId}/setting/wlans/{wlanId}/ssids

	// call patch
	// PATCH /{omadacId}/api/v2/sites/{siteId}/setting/wlans/{wlanId}/ssids/{ssidId}

	c.writeResponse(w, http.StatusOK, true, nil)

}

func (c *controller) writeResponse(w http.ResponseWriter, statusCode int, updated bool, errorMessage *string) {
	var nw NetworksResponse
	if errorMessage != nil {
		nw.Error = &NetworksResponseError{
			Code:    statusCode,
			Message: *errorMessage,
		}
		glog.Error(nw.Error.Message)
	}

	if updated {
		nw.Data = &NetworksResponseSuccess{
			Updated: true,
		}
	}

	data, err := json.Marshal(nw)
	if err != nil {
		glog.Errorf("Cannot marshal json response: %s", err)
		http.Error(w, "Cannot marshal json response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		glog.Errorf("Cannot write a response. %s", err)
		http.Error(w, "Cannot write a response.", http.StatusUnprocessableEntity)
	}
}
