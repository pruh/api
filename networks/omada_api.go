package networks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/pruh/api/config"
	apihttp "github.com/pruh/api/http"
)

// OmadaApi to interact with an Omada API web servies
type omadaApi struct {
	config     *config.Configuration
	httpClient apihttp.Client
}

type OmadaApi interface {
	GetControllerId() (*OmadaResponse, error)
	Login(omadaControllerId *string) (*OmadaResponse, error)
}

// NewOmadaApi creates new omada api
func NewOmadaApi(config *config.Configuration, httpClient apihttp.Client) OmadaApi {
	return &omadaApi{
		config:     config,
		httpClient: httpClient,
	}
}

func (oa *omadaApi) GetControllerId() (*OmadaResponse, error) {
	req, err := http.NewRequest(http.MethodGet, *oa.config.OmadaUrl+"/api/info", nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	glog.Infof("sending GetControllerId request")

	resp, err := oa.httpClient.Do(req)
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

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada controller id: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) Login(omadaControllerId *string) (*OmadaResponse, error) {
	l := LoginData{
		Username: *oa.config.OmadaUsername,
		Password: *oa.config.OmadaPassword,
	}

	jsonStr, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/api/v2/login", *oa.config.OmadaUrl, *omadaControllerId)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	glog.Infof("sending login request")

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error logging in: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing login: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
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
