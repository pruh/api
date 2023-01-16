package networks

// nested within sbserver response
type OmadaResponse struct {
	ErrorCode int     `json:"errorCode"`
	Msg       *string `json:"msg,omitempty"`
	Result    *Result `json:"result,omitempty"`
}

type Result struct {
	OmadacId  *string `json:"omadacId,omitempty"`
	Token     *string `json:"token,omitempty"`
	ProfileId *string `json:"profileId,omitempty"`
	Data      *[]Data `json:"data,omitempty"`
}

type Data struct {
	Id       *string     `json:"id,omitempty"`
	Name     *string     `json:"name,omitempty"`
	DayMode  *int        `json:"dayMode,omitempty"`
	TimeList *[]TimeList `json:"timeList,omitempty"`
	DayMon   *bool       `json:"dayMon,omitempty"`
	DayTue   *bool       `json:"dayTue,omitempty"`
	DayWed   *bool       `json:"dayWed,omitempty"`
	DayThu   *bool       `json:"dayThu,omitempty"`
	DayFri   *bool       `json:"dayFri,omitempty"`
	DaySat   *bool       `json:"daySat,omitempty"`
	DaySun   *bool       `json:"daySun,omitempty"`
}

type TimeList struct {
	DayType    *int `json:"dayType,omitempty"`
	StartTimeH *int `json:"startTimeH,omitempty"`
	StartTimeM *int `json:"startTimeM,omitempty"`
	EndTimeH   *int `json:"endTimeH,omitempty"`
	EndTimeM   *int `json:"endTimeM,omitempty"`
}

type OmadaLoginData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type OmadaSsidUpdateData struct {
	WlanScheduleEnable *bool   `json:"wlanScheduleEnable,omitempty"`
	Action             *int    `json:"action,omitempty"`
	ScheduleId         *string `json:"scheduleId,omitempty"`
}

type OmadaTimeRangeData struct {
	Name     *string     `json:"name,omitempty"`
	DayMode  *int        `json:"dayMode,omitempty"`
	TimeList *[]TimeList `json:"timeList,omitempty"`
}

type NetworksResponse struct {
	Data  *NetworksResponseSuccess `json:"data,omitempty"`
	Error *NetworksResponseError   `json:"error,omitempty"`
}

type NetworksResponseSuccess struct {
	Updated bool `json:"updated,omitempty"`
}

type NetworksResponseError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewStr(str string) *string {
	return &str
}

func NewInt(num int) *int {
	return &num
}

func NewBool(b bool) *bool {
	return &b
}
