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
		c.writeResponse(w, http.StatusBadRequest, false,
			NewStr("ssid is missing in the request parameters"), nil, nil)
		return
	}

	//controller id
	omadaIdResp, err := c.repository.GetControllerId()
	if err != nil || omadaIdResp == nil ||
		omadaIdResp.ErrorCode != 0 || omadaIdResp.Result.OmadacId == nil {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada controller query error"), omadaIdResp, err)
		return
	}

	glog.Infof("Omada controller id %s", *omadaIdResp.Result.OmadacId)

	omadaLoginResp, cookies, err := c.repository.Login(omadaIdResp.Result.OmadacId)
	if err != nil || omadaLoginResp == nil || omadaLoginResp.ErrorCode != 0 || cookies == nil ||
		omadaLoginResp.Result == nil || omadaLoginResp.Result.Token == nil {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada login query error"), omadaLoginResp, err)
		return
	}

	glog.Infof("Omada login token %s", *omadaLoginResp.Result.Token)

	omadaSitesResp, err := c.repository.GetSites(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token)
	if err != nil || omadaSitesResp == nil ||
		omadaSitesResp.ErrorCode != 0 || omadaSitesResp.Result == nil ||
		omadaSitesResp.Result.Data == nil || len(*omadaSitesResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada sites query error"), omadaSitesResp, err)
		return
	}

	glog.Infof("Omada sites %+v", (*omadaSitesResp.Result.Data)[0])

	omadaWlansResp, err := c.repository.GetWlans(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaWlansResp == nil ||
		omadaWlansResp.ErrorCode != 0 || omadaWlansResp.Result == nil ||
		omadaWlansResp.Result.Data == nil || len(*omadaWlansResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada wlans query error"), omadaWlansResp, err)
		return
	}

	glog.Infof("Omada wlans %+v", (*omadaWlansResp.Result.Data)[0])

	omadaSsidsResp, err := c.repository.GetSsids(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id, (*omadaWlansResp.Result.Data)[0].Id)
	if err != nil || omadaSsidsResp == nil ||
		omadaSsidsResp.ErrorCode != 0 || omadaSsidsResp.Result == nil ||
		omadaSsidsResp.Result.Data == nil || len(*omadaSsidsResp.Result.Data) == 0 {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada ssids query error"), omadaSsidsResp, err)
		return
	}

	glog.Infof("Omada ssids %+v", (*omadaSsidsResp.Result.Data)[0])

	var ssidId *string
	for _, ssidData := range *omadaSsidsResp.Result.Data {
		if *ssidData.Name == ssid {
			ssidId = ssidData.Id
			break
		}
	}

	if ssidId == nil {
		c.writeResponse(w, http.StatusNotFound, false,
			NewStr("ssid with given name not found in configured networks"), nil, nil)
		return
	}

	glog.Infof("Omada ssid id %s", *ssidId)

	omadaTimeRangesResp, err := c.repository.GetTimeRanges(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id)
	if err != nil || omadaTimeRangesResp == nil || omadaTimeRangesResp.ErrorCode != 0 {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Omada time ranges query error"), omadaSitesResp, err)
		return
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
			omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id,
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

		if err != nil || omadaTrCreateResp == nil ||
			omadaTrCreateResp.ErrorCode != 0 || omadaTrCreateResp.Result.ProfileId == nil {
			c.writeResponse(w, http.StatusBadGateway, false,
				NewStr("Can not create time range"), omadaTrCreateResp, err)
			return
		}

		scheduleId = omadaTrCreateResp.Result.ProfileId
	}

	omadaUpdateSsidResp, err := c.repository.UpdateSsid(omadaIdResp.Result.OmadacId, cookies,
		omadaLoginResp.Result.Token, (*omadaSitesResp.Result.Data)[0].Id,
		(*omadaWlansResp.Result.Data)[0].Id, &ssid, &OmadaSsidUpdateData{
			WlanScheduleEnable: NewBool(true),
			Action:             NewInt(0),
			ScheduleId:         scheduleId,
		})
	if err != nil || omadaUpdateSsidResp.ErrorCode != 0 {
		c.writeResponse(w, http.StatusBadGateway, false,
			NewStr("Can not update ssid"), omadaUpdateSsidResp, err)
		return
	}

	c.writeResponse(w, http.StatusOK, true, nil, nil, nil)

}

func (c *controller) writeResponse(w http.ResponseWriter, statusCode int, updated bool,
	errorMessage *string, upstreamResp *OmadaResponse, upstreamErr error) {
	var nw NetworksResponse
	if errorMessage != nil {
		var fmtErrorMessage *string
		if upstreamResp != nil {
			fmtErrorMessage = NewStr(fmt.Sprintf("%s: %s", *errorMessage, *upstreamResp.Msg))
		} else if upstreamErr != nil {
			fmtErrorMessage = NewStr(fmt.Sprintf("%s: %s", *errorMessage, upstreamErr))
		} else {
			fmtErrorMessage = errorMessage
		}

		nw.Error = &NetworksResponseError{
			Code:    statusCode,
			Message: *fmtErrorMessage,
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
