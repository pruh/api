package networks

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/pruh/api/config"
	apihttp "github.com/pruh/api/http"
)

// OmadaApi to interact with an Omada API web servies
type omadaApi struct {
	Config     *config.Configuration
	HTTPClient apihttp.Client
}

type OmadaApi interface {
	GetControllerId() (*ControllerIdResponse, error)
}

// NewOmadaApi creates new omada api
func NewOmadaApi(config *config.Configuration, httpClient apihttp.Client) OmadaApi {
	return &omadaApi{
		Config:     config,
		HTTPClient: httpClient,
	}
}

func (oa *omadaApi) GetControllerId() (*ControllerIdResponse, error) {
	req, err := http.NewRequest(http.MethodGet, *oa.Config.OmadaUrl+"/api/info", nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	glog.Infof("sending GetControllerId request")

	resp, err := oa.HTTPClient.Do(req)

	if err != nil {
		glog.Errorf("Error querying omada controller id: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada response: %s", err)
		return nil, err
	}

	controllerIdResponse := &ControllerIdResponse{}
	err = json.Unmarshal(body, &controllerIdResponse)
	if err != nil {
		glog.Errorf("Error parsing omada controller id: %s", err)
		return nil, err
	}

	return controllerIdResponse, nil
}

// func (oa *OmadaApi) Login() string {
// 	return nil
// }

// func (oa *OmadaApi) GetSiteId() string {
// 	return nil
// }

// func (oa *OmadaApi) GetWlanId(siteId *string) string {
// 	return nil
// }

// func (oa *OmadaApi) GetSsidId(siteId *string, wlanId *string) string {
// 	return nil
// }
