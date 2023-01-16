package networks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	GetSites(omadaControllerId *string, loginToken *string) (*OmadaResponse, error)
	GetWlans(omadaControllerId *string, loginToken *string,
		siteId *string) (*OmadaResponse, error)
	GetSsids(omadaControllerId *string, loginToken *string,
		siteId *string, wlanId *string) (*OmadaResponse, error)
	UpdateSsid(omadaControllerId *string, loginToken *string,
		siteId *string, wlanId *string, ssidId *string,
		ssidUpdateData *OmadaSsidUpdateData) (*OmadaResponse, error)
	GetTimeRanges(omadaControllerId *string, loginToken *string,
		siteId *string) (*OmadaResponse, error)
	CreateTimeRange(omadaControllerId *string, loginToken *string,
		siteId *string, trData *Data) (*OmadaResponse, error)
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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
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
	l := OmadaLoginData{
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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
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

func (oa *omadaApi) GetSites(omadaControllerId *string, loginToken *string) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites?currentPageSize=1&currentPage=1", *oa.config.OmadaUrl, *omadaControllerId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	glog.Infof("sending GetSites request")

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error querying omada sites: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada sites response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada sites: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) GetWlans(omadaControllerId *string, loginToken *string, siteId *string) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans", *oa.config.OmadaUrl, *omadaControllerId, *siteId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	glog.Infof("sending GetWlans request")

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error querying omada wlans: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada wlans response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada wlans: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) GetSsids(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans/%s/ssids",
		*oa.config.OmadaUrl, *omadaControllerId, *siteId, *wlanId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	glog.Infof("sending GetSsids request")

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error querying omada ssids: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada ssids response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada ssids: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) UpdateSsid(omadaControllerId *string, loginToken *string,
	siteId *string, wlanId *string, ssidId *string,
	ssidUpdateData *OmadaSsidUpdateData) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/wlans/%s/ssids/%s",
		*oa.config.OmadaUrl, *omadaControllerId, *siteId, *wlanId, *ssidId)

	glog.Infof("updating ssid %s with %+v", *ssidId, ssidUpdateData)

	jsonStr, err := json.Marshal(ssidUpdateData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error updating omada ssids: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada ssids response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada ssids: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) GetTimeRanges(omadaControllerId *string, loginToken *string,
	siteId *string) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/profiles/timeranges",
		*oa.config.OmadaUrl, *omadaControllerId, *siteId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	glog.Infof("sending GetTimeRanges request")

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error querying omada time ranges: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada time ranges response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada time ranges: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}

func (oa *omadaApi) CreateTimeRange(omadaControllerId *string, loginToken *string,
	siteId *string, trData *Data) (*OmadaResponse, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/profiles/timeranges",
		*oa.config.OmadaUrl, *omadaControllerId, *siteId)

	glog.Infof("creating new time range %+v", *trData)

	jsonStr, err := json.Marshal(trData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		glog.Errorf("Failed to create HTTP request: %s", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Csrf-token", *loginToken)

	resp, err := oa.httpClient.Do(req)
	if err != nil {
		glog.Errorf("Error creating time range: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading omada time range create response: %s", err)
		return nil, err
	}

	var omadaResponse OmadaResponse
	err = json.Unmarshal(body, &omadaResponse)
	if err != nil {
		glog.Errorf("Error parsing omada time range create response: %s", err)
		return nil, err
	}

	return &omadaResponse, nil
}
