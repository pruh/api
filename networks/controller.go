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
	Config     *config.Configuration
	Repository Repository
}

type Controller interface {
	UpdateWifis(w http.ResponseWriter, r *http.Request)
}

// NewController creates new networks controller
func NewController(config *config.Configuration) Controller {
	omadaApi := NewOmadaApi(config, apihttp.NewHTTPClient())
	return NewControllerWithParams(config, omadaApi)
}

// Creates a new controller with additional dependencies for tests
func NewControllerWithParams(config *config.Configuration, omadaApi OmadaApi) Controller {
	return &controller{
		Config:     config,
		Repository: NewRepository(omadaApi),
	}
}

func (c *controller) UpdateWifis(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		errorMessage := "ssid is missing in the request parameters"
		c.writeResponse(w, http.StatusBadRequest, false, &errorMessage)
		return
	}

	//controller id
	omadaId, err := c.Repository.OmadaApi.GetControllerId()
	if err != nil || omadaId.ErrorCode != 0 || len(omadaId.Result.OmadacId) == 0 {
		errorMessage := fmt.Sprintf("Omada Controller Id Query Error: %+v", err)
		c.writeResponse(w, http.StatusBadGateway, false, &errorMessage)
		return
	}

	glog.Infof("Omada controller id %s", omadaId.Result.OmadacId)

	// login

	// obtain login token

	// get site id

	// get wlan id

	// get ssidid

	// call patch

	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// c.Config.OmadaUrl

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
