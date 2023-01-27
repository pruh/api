package networks

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

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
