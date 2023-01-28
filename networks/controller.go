package networks

import (
	"encoding/json"
	"errors"
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
	ufc        UrlFilterController
	repository Repository
}

type Controller interface {
	GetWifi(w http.ResponseWriter, r *http.Request)
	UpdateWifi(w http.ResponseWriter, r *http.Request)
}

// Creates new networks controller
func NewController(config *config.Configuration) Controller {
	omadaApi := NewOmadaApi(config, apihttp.NewHTTPClient())
	r := NewRepository(omadaApi)
	ufc := NewUrlFilterController(r)
	return NewControllerWithParams(config, ufc, r)
}

// Creates a new controller with additional dependencies for tests
func NewControllerWithParams(config *config.Configuration,
	ufc UrlFilterController, r Repository) Controller {
	return &controller{
		config:     config,
		ufc:        ufc,
		repository: r,
	}
}

func (c *controller) GetWifi(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, ssid, nil,
			errors.New("ssid is missing in the request parameters"))
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp == nil ||
		omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaIdResp,
			fmt.Errorf("omada controller query error %v", err))
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, cookies, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaLoginResp == nil || omadaLoginResp.ErrorCode != 0 || cookies == nil ||
		omadaLoginResp.Result == nil || omadaLoginResp.Result.Token == nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaLoginResp,
			fmt.Errorf("omada login query error %v", err))
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token)
	if err != nil || omadaSitesResp == nil ||
		omadaSitesResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaSitesResp,
			fmt.Errorf("omada sites query error %v", err))
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	omadaWlansResp, err := c.repository.GetWlans(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaWlansResp == nil ||
		omadaWlansResp.ErrorCode != 0 || omadaWlansResp.Result == nil ||
		omadaWlansResp.Result.Data == nil || len(*omadaWlansResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaWlansResp,
			fmt.Errorf("omada wlans query error %v", err))
		return
	}

	glog.Infof("Omada wlans %+v", (*omadaWlansResp.Result.Data)[0])

	omadaSsidsResp, err := c.repository.GetSsids(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, (*omadaWlansResp.Result.Data)[0].Id)
	if err != nil || omadaSsidsResp == nil ||
		omadaSsidsResp.ErrorCode != 0 || omadaSsidsResp.Result == nil ||
		omadaSsidsResp.Result.Data == nil || len(*omadaSsidsResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaSsidsResp,
			fmt.Errorf("omada ssids query error %v", err))
		return
	}

	glog.Infof("Omada ssids %+v", (*omadaSsidsResp.Result.Data)[0])

	var ssidData *Data
	for _, sd := range *omadaSsidsResp.Result.Data {
		if *sd.Name == ssid {
			ssidData = &sd
			break
		}
	}

	if ssidData == nil {
		writeErrorResponse(w, http.StatusNotFound, ssid, nil,
			errors.New("ssid with given name not found in configured networks"))
		return
	}

	glog.Infof("Omada ssid id %s", *ssidData.Id)

	urlFilters, err := c.ufc.QueryUrlFilters(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, ssidData)
	if err != nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaSsidsResp,
			fmt.Errorf("failure querying url filters %v", err))
		return
	}

	writeSuccessResponse(w, ssid, ssidData, urlFilters, nil)
}

func (c *controller) UpdateWifi(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, ssid, nil,
			errors.New("ssid is missing in the request parameters"))
		return
	}

	var ssidRequest NetworksSsidRequest
	err := json.NewDecoder(r.Body).Decode(&ssidRequest)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, ssid, nil,
			fmt.Errorf("request is malformed %v", err))
		return
	}

	if ssidRequest.RadioOn == nil &&
		ssidRequest.UploadLimit == nil && ssidRequest.DownloadLimit == nil &&
		ssidRequest.UrlFilters == nil {
		writeSuccessResponse(w, ssid, nil, nil, NewBool(false))
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp == nil ||
		omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaIdResp,
			fmt.Errorf("omada controller query error %v", err))
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, cookies, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaLoginResp == nil || omadaLoginResp.ErrorCode != 0 || cookies == nil ||
		omadaLoginResp.Result == nil || omadaLoginResp.Result.Token == nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaLoginResp,
			fmt.Errorf("omada login query error %v", err))
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token)
	if err != nil || omadaSitesResp == nil ||
		omadaSitesResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaSitesResp,
			fmt.Errorf("omada sites query error %v", err))
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	omadaWlansResp, err := c.repository.GetWlans(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaWlansResp == nil ||
		omadaWlansResp.ErrorCode != 0 || omadaWlansResp.Result == nil ||
		omadaWlansResp.Result.Data == nil || len(*omadaWlansResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaWlansResp,
			fmt.Errorf("omada wlans query error %v", err))
		return
	}

	glog.Infof("Omada wlans %+v", (*omadaWlansResp.Result.Data)[0])

	omadaSsidsResp, err := c.repository.GetSsids(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, (*omadaWlansResp.Result.Data)[0].Id)
	if err != nil || omadaSsidsResp == nil ||
		omadaSsidsResp.ErrorCode != 0 || omadaSsidsResp.Result == nil ||
		omadaSsidsResp.Result.Data == nil || len(*omadaSsidsResp.Result.Data) == 0 {
		writeErrorResponse(w, http.StatusBadGateway, ssid, omadaSsidsResp,
			fmt.Errorf("omada ssids query error %v", err))
		return
	}

	glog.Infof("Omada ssids %+v", (*omadaSsidsResp.Result.Data)[0])

	var ssidData *Data
	for _, sd := range *omadaSsidsResp.Result.Data {
		if *sd.Name == ssid {
			ssidData = &sd
			break
		}
	}

	if ssidData == nil {
		writeErrorResponse(w, http.StatusNotFound, ssid, nil,
			errors.New("ssid with given name not found in configured networks"))
		return
	}

	glog.Infof("Omada ssid id %s", *ssidData.Id)

	needtoUpdate := false
	if needRadioStateUpdate(ssidData.WlanScheduleEnable, ssidRequest.RadioOn) {
		glog.Info("need to update radio state")

		needtoUpdate = true

		ssidData.WlanScheduleEnable = NewBool(!*ssidRequest.RadioOn)
		if *ssidData.WlanScheduleEnable {
			glog.Infof("Looking for time range for ssid %s", *ssidData.Id)

			scheduleId, err := c.getTimeRange(omadaIdResp, cookies,
				omadaLoginResp.Result.Token, omadaSitesResp)
			if err != nil {
				writeErrorResponse(w, http.StatusBadGateway, ssid, nil, err)
				return
			}

			ssidData.Action = NewInt(0)
			ssidData.ScheduleId = scheduleId
		}
	}

	if !IsSpeedLimitEqual(ssidData.RateLimit.DownLimitEnable, ssidData.RateLimit.DownLimit,
		ssidData.RateLimit.DownLimitType, ssidRequest.DownloadLimit) {

		glog.Info("need to update download limit")

		needtoUpdate = true

		ssidData.RateLimit.DownLimitEnable = NewBool(*ssidRequest.DownloadLimit > 0)
		if *ssidRequest.DownloadLimit > 0 {
			ssidData.RateLimit.DownLimit = NewInt(*ssidRequest.DownloadLimit)
		} else {
			ssidData.RateLimit.DownLimit = NewInt(0)
		}
		ssidData.RateLimit.DownLimitType = NewInt(0)

		ssidData.RateLimit.RateLimitId = nil
	}

	if !IsSpeedLimitEqual(ssidData.RateLimit.UpLimitEnable, ssidData.RateLimit.UpLimit,
		ssidData.RateLimit.UpLimitType, ssidRequest.UploadLimit) {

		glog.Info("need to update upload limit")

		needtoUpdate = true

		ssidData.RateLimit.UpLimitEnable = NewBool(*ssidRequest.UploadLimit > 0)
		if *ssidRequest.UploadLimit > 0 {
			ssidData.RateLimit.UpLimit = NewInt(*ssidRequest.UploadLimit)
		} else {
			ssidData.RateLimit.UpLimit = NewInt(0)
		}
		ssidData.RateLimit.UpLimitType = NewInt(0)

		ssidData.RateLimit.RateLimitId = nil
	}

	if needtoUpdate {
		omadaUpdateSsidResp, err := c.repository.UpdateSsid(omadaIdResp.Result.OmadacId, cookies,
			omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id,
			(*omadaWlansResp.Result.Data)[0].Id, ssidData)

		if err != nil || omadaUpdateSsidResp == nil || omadaUpdateSsidResp.ErrorCode != 0 {
			writeErrorResponse(w, http.StatusBadGateway, ssid, omadaUpdateSsidResp,
				fmt.Errorf("can not update ssid %v", err))
			return
		}
	}

	uf, updated, err := c.ufc.MaybeUpdateUrlFilters(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, ssidData, ssidRequest.UrlFilters)
	if err != nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, nil,
			fmt.Errorf("can not update url filters %v", err))
		return
	}

	needtoUpdate = needtoUpdate || *updated

	writeSuccessResponse(w, ssid, ssidData, uf, NewBool(needtoUpdate))
}

func needRadioStateUpdate(wlanScheduleEnable *bool, requestRadioOn *bool) bool {
	// radio state needs update if request current radio state is equal to schedule enabled state
	if requestRadioOn == nil {
		return false
	}
	return *requestRadioOn == *wlanScheduleEnable
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, ssid string,
	upstreamResp *OmadaResponse, upstreamErr error) {

	var nw NetworksResponse
	if upstreamErr != nil {
		var fmtErrorMessage *string
		if upstreamResp != nil {
			fmtErrorMessage = NewStr(fmt.Sprintf("%s: %v", *upstreamResp.Msg, upstreamErr))
		} else {
			fmtErrorMessage = NewStr(fmt.Sprintf("%v", upstreamErr))
		}

		nw.ErrorMessage = fmtErrorMessage
	}

	nw.Ssid = &ssid

	writeResponse(w, statusCode, &nw)
}

func writeSuccessResponse(w http.ResponseWriter, ssid string, ssidData *Data,
	urlFilters *[]UrlFilter, updated *bool) {

	if ssidData == nil {
		writeResponse(w, http.StatusOK, &NetworksResponse{Ssid: &ssid, Updated: updated})
		return
	}

	upRate, err := ssidData.ToUploadRateLimit()
	if err != nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, nil,
			fmt.Errorf("omada ssid upload rate limit error %v", err))
		return
	}
	downRate, err := ssidData.ToDownloadRateLimit()
	if err != nil {
		writeErrorResponse(w, http.StatusBadGateway, ssid, nil,
			fmt.Errorf("omada ssid download rate limit error %v", err))
		return
	}
	var nw NetworksResponse

	nw.Ssid = &ssid
	nw.RadioOn = NewBool(!*ssidData.WlanScheduleEnable)
	nw.UploadLimit = upRate
	nw.DownloadLimit = downRate
	nw.UrlFilters = urlFilters
	nw.Updated = updated

	writeResponse(w, http.StatusOK, &nw)
}

func writeResponse(w http.ResponseWriter, statusCode int, nw *NetworksResponse) {
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
