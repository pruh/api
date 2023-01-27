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
	repository Repository
}

type Controller interface {
	GetWifi(w http.ResponseWriter, r *http.Request)
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

func (c *controller) GetWifi(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		c.writeResponse(w, http.StatusBadRequest, nil, nil, nil, nil, nil, nil,
			errors.New("ssid is missing in the request parameters"))
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp == nil ||
		omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaIdResp,
			fmt.Errorf("omada controller query error %v", err))
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, cookies, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaLoginResp == nil || omadaLoginResp.ErrorCode != 0 || cookies == nil ||
		omadaLoginResp.Result == nil || omadaLoginResp.Result.Token == nil {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaLoginResp,
			fmt.Errorf("omada login query error %v", err))
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token)
	if err != nil || omadaSitesResp == nil ||
		omadaSitesResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaSitesResp,
			fmt.Errorf("omada sites query error %v", err))
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	omadaWlansResp, err := c.repository.GetWlans(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaWlansResp == nil ||
		omadaWlansResp.ErrorCode != 0 || omadaWlansResp.Result == nil ||
		omadaWlansResp.Result.Data == nil || len(*omadaWlansResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaWlansResp,
			fmt.Errorf("omada wlans query error %v", err))
		return
	}

	glog.Infof("Omada wlans %+v", (*omadaWlansResp.Result.Data)[0])

	omadaSsidsResp, err := c.repository.GetSsids(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, (*omadaWlansResp.Result.Data)[0].Id)
	if err != nil || omadaSsidsResp == nil ||
		omadaSsidsResp.ErrorCode != 0 || omadaSsidsResp.Result == nil ||
		omadaSsidsResp.Result.Data == nil || len(*omadaSsidsResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaSsidsResp,
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
		c.writeResponse(w, http.StatusNotFound, nil, nil, nil, nil, nil, nil,
			errors.New("ssid with given name not found in configured networks"))
		return
	}

	glog.Infof("Omada ssid id %s", *ssidData.Id)

	// todo add converter
	u, err := ssidData.ToUploadRateLimit()
	d, err := ssidData.ToDownloadRateLimit()
	c.writeResponse(w, http.StatusOK, ssidData.Name,
		NewBool(!*ssidData.WlanScheduleEnable), u, d,
		nil, nil, nil)
}

func (c *controller) UpdateWifi(w http.ResponseWriter, r *http.Request) {
	// verify ssid is present
	params := mux.Vars(r)
	ssid := params["ssid"]
	if len(ssid) == 0 {
		c.writeResponse(w, http.StatusBadRequest, nil, nil, nil, nil, nil, nil,
			errors.New("ssid is missing in the request parameters"))
		return
	}

	var ssidRequest NetworksSsidRequest
	err := json.NewDecoder(r.Body).Decode(&ssidRequest)
	if err != nil {
		c.writeResponse(w, http.StatusBadRequest, nil, nil, nil, nil, nil, nil,
			fmt.Errorf("request is malformed %v", err))
		return
	}

	if ssidRequest.RadioOn == nil && ssidRequest.UploadLimit == nil && ssidRequest.DownloadLimit == nil {
		c.writeResponse(w, http.StatusOK, &ssid, nil, nil, nil,
			NewBool(false), nil, nil)
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp == nil ||
		omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaIdResp,
			fmt.Errorf("omada controller query error %v", err))
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, cookies, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaLoginResp == nil || omadaLoginResp.ErrorCode != 0 || cookies == nil ||
		omadaLoginResp.Result == nil || omadaLoginResp.Result.Token == nil {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaLoginResp,
			fmt.Errorf("omada login query error %v", err))
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token)
	if err != nil || omadaSitesResp == nil ||
		omadaSitesResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaSitesResp,
			fmt.Errorf("omada sites query error %v", err))
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	omadaWlansResp, err := c.repository.GetWlans(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaWlansResp == nil ||
		omadaWlansResp.ErrorCode != 0 || omadaWlansResp.Result == nil ||
		omadaWlansResp.Result.Data == nil || len(*omadaWlansResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaWlansResp,
			fmt.Errorf("omada wlans query error %v", err))
		return
	}

	glog.Infof("Omada wlans %+v", (*omadaWlansResp.Result.Data)[0])

	omadaSsidsResp, err := c.repository.GetSsids(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, (*omadaWlansResp.Result.Data)[0].Id)
	if err != nil || omadaSsidsResp == nil ||
		omadaSsidsResp.ErrorCode != 0 || omadaSsidsResp.Result == nil ||
		omadaSsidsResp.Result.Data == nil || len(*omadaSsidsResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaSsidsResp,
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
		c.writeResponse(w, http.StatusNotFound, nil, nil, nil, nil, nil, nil,
			errors.New("ssid with given name not found in configured networks"))
		return
	}

	glog.Infof("Omada ssid id %s", *ssidData.Id)

	if c.isRadioStateEqual(ssidData.WlanScheduleEnable, ssidRequest.RadioOn) &&
		c.isSpeedLimitEqual(ssidData.RateLimit.DownLimitEnable, ssidData.RateLimit.DownLimit,
			ssidData.RateLimit.DownLimitType, ssidRequest.DownloadLimit) &&
		c.isSpeedLimitEqual(ssidData.RateLimit.UpLimitEnable, ssidData.RateLimit.UpLimit,
			ssidData.RateLimit.UpLimitType, ssidRequest.UploadLimit) {

		glog.Info("no need to update ssid")
		// TODO update
		c.writeResponse(w, http.StatusOK, ssidData.Name,
			NewBool(!*ssidData.WlanScheduleEnable), NewInt(0), NewInt(0),
			NewBool(false), nil, nil)
		return
	}

	// TODO only if not null
	ssidData.WlanScheduleEnable = NewBool(!*ssidRequest.RadioOn)
	if *ssidData.WlanScheduleEnable {
		glog.Infof("Looking for time range for ssid %s", *ssidData.Id)

		scheduleId, err := c.getTimeRange(omadaIdResp, cookies, omadaLoginResp.Result.Token, omadaSitesResp)
		if err != nil {
			c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, nil, err)
			return
		}

		ssidData.Action = NewInt(0)
		ssidData.ScheduleId = scheduleId
	}

	omadaUpdateSsidResp, err := c.repository.UpdateSsid(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id,
		(*omadaWlansResp.Result.Data)[0].Id, ssidData)

	if err != nil || omadaUpdateSsidResp == nil || omadaUpdateSsidResp.ErrorCode != 0 {
		c.writeResponse(w, http.StatusBadGateway, nil, nil, nil, nil, nil, omadaUpdateSsidResp,
			fmt.Errorf("can not update ssid %v", err))
		return
	}

	// TODO update
	c.writeResponse(w, http.StatusOK, ssidData.Name,
		NewBool(!*ssidData.WlanScheduleEnable), NewInt(0), NewInt(0),
		NewBool(true), nil, nil)
}

func (c *controller) isRadioStateEqual(wlanScheduleEnable *bool, requestRadioOn *bool) bool {
	// state is equal if request radio state is NOT equal schedule enabled state
	return *requestRadioOn == !*wlanScheduleEnable
}

func (c *controller) isSpeedLimitEqual(
	speedLimitEnable *bool,
	speedLimit *int,
	speedLimitType *int,
	requestSpeedLimit *int) bool {
	if requestSpeedLimit == nil {
		// no speed limit in request
		return true
	}

	if *requestSpeedLimit < 1 {
		// request to set no speed limit
		// speed is equal if speed limit not set
		return !*speedLimitEnable
	}

	// request to set speed limit

	if !*speedLimitEnable {
		// speed is NOT equal if speed limit is NOT enabled
		return false
	}

	// speed is equal if speed limit speed is the same
	speedLimitKbps := *speedLimit * (1024 ^ *speedLimitType)
	return speedLimitKbps == *requestSpeedLimit
}

func (c *controller) getTimeRange(omadaIdResp *OmadaResponse, cookies []*http.Cookie,
	token *string, omadaSitesResp *OmadaResponse) (*string, error) {
	omadaTimeRangesResp, err := c.repository.GetTimeRanges(omadaIdResp.Result.OmadacId, cookies,
		token, (*omadaSitesResp.Result.Data)[0].Id)

	if err != nil || omadaTimeRangesResp == nil {
		return nil, fmt.Errorf("omada time ranges query error %v", err)
	} else if omadaTimeRangesResp.ErrorCode != 0 {
		return nil, fmt.Errorf("omada time ranges query error %s", *omadaTimeRangesResp.Msg)
	}

	var timeRangeData *Data
	if omadaTimeRangesResp.Result.Data != nil {
		for _, tr := range *omadaTimeRangesResp.Result.Data {
			if tr.DayMode != nil && *tr.DayMode == 0 && tr.TimeList != nil && len(*tr.TimeList) > 0 {
				for _, tl := range *tr.TimeList {
					if *tl.StartTimeH == 0 && *tl.StartTimeM == 0 && *tl.EndTimeH == 24 && *tl.EndTimeM == 0 {
						timeRangeData = &tr
						break
					}
				}
			}
		}
	}

	var scheduleId *string
	if timeRangeData != nil {
		glog.Infof("time range already exists: %+v", *timeRangeData)
		scheduleId = timeRangeData.Id
	} else {
		glog.Info("time range not available, creating one")

		omadaTrCreateResp, err := c.repository.CreateTimeRange(omadaIdResp.Result.OmadacId, cookies,
			token, (*omadaSitesResp.Result.Data)[0].Id,
			&Data{
				Name:    NewStr("Night and Day"),
				DayMode: NewInt(0),
				DayMon:  NewBool(true),
				DayTue:  NewBool(true),
				DayWed:  NewBool(true),
				DayThu:  NewBool(true),
				DayFri:  NewBool(true),
				DaySat:  NewBool(true),
				DaySun:  NewBool(true),
				TimeList: &[]TimeList{
					{
						DayType:    NewInt(0),
						StartTimeH: NewInt(0),
						StartTimeM: NewInt(0),
						EndTimeH:   NewInt(24),
						EndTimeM:   NewInt(0),
					},
				},
			})

		if err != nil || omadaTrCreateResp == nil {
			return nil, fmt.Errorf("can not create time range %v", err)
		} else if omadaTrCreateResp.ErrorCode != 0 || omadaTrCreateResp.Result.ProfileId == nil {
			return nil, fmt.Errorf("can not create time range %s", *omadaTrCreateResp.Msg)
		}

		scheduleId = omadaTrCreateResp.Result.ProfileId
	}

	return scheduleId, nil
}

func (c *controller) writeResponse(w http.ResponseWriter, statusCode int, ssid *string,
	radioOn *bool, uploadLimit *int, downloadLimit *int, updated *bool,
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

		glog.Error(*nw.ErrorMessage)
	}

	nw.Ssid = ssid
	nw.RadioOn = radioOn
	nw.UploadLimit = uploadLimit
	nw.DownloadLimit = downloadLimit
	nw.Updated = updated

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
