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
	// common params
	Id     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	SiteId *string `json:"siteid,omitempty"`
	Policy *int    `json:"policy,omitempty"`

	// ssid params
	Band               *int               `json:"band,omitempty"`
	WlanId             *string            `json:"wlanid,omitempty"`
	VlanEnable         *bool              `json:"vlanEnable,omitempty"`
	VlanId             *string            `json:"vlanId,omitempty"`
	Broadcast          *bool              `json:"broadcast,omitempty"`
	Security           *int               `json:"security,omitempty"`
	GuestNetEnable     *bool              `json:"guestNetEnable,omitempty"`
	WlanScheduleEnable *bool              `json:"wlanScheduleEnable,omitempty"`
	Action             *int               `json:"action,omitempty"`
	ScheduleId         *string            `json:"scheduleId,omitempty"`
	MacFilterEnable    *bool              `json:"macFilterEnable,omitempty"`
	MacFilterId        *string            `json:"macFilterId,omitempty"`
	RateLimit          *RateLimit         `json:"rateLimit,omitempty"`
	PskSetting         *PskSetting        `json:"pskSetting,omitempty"`
	WpaSetting         *WpaSetting        `json:"wpaSetting,omitempty"`
	RateAndBeaconCtrl  *RateAndBeaconCtrl `json:"rateAndBeaconCtrl,omitempty"`

	// time range params
	DayMode  *int        `json:"dayMode,omitempty"`
	TimeList *[]TimeList `json:"timeList,omitempty"`
	DayMon   *bool       `json:"dayMon,omitempty"`
	DayTue   *bool       `json:"dayTue,omitempty"`
	DayWed   *bool       `json:"dayWed,omitempty"`
	DayThu   *bool       `json:"dayThu,omitempty"`
	DayFri   *bool       `json:"dayFri,omitempty"`
	DaySat   *bool       `json:"daySat,omitempty"`
	DaySun   *bool       `json:"daySun,omitempty"`

	// url filtering
	Type       *string   `json:"type,omitempty"`
	EntryId    *int      `json:"entryId,omitempty"`
	Status     *bool     `json:"status,omitempty"`
	SourceType *int      `json:"sourceType,omitempty"`
	SourceIds  *[]string `json:"sourceIds,omitempty"`
	Urls       *[]string `json:"urls,omitempty"`
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

type PskSetting struct {
	SecurityKey       *string `json:"securityKey,omitempty"`
	VersionPsk        *int    `json:"versionPsk,omitempty"`
	EncryptionPsk     *int    `json:"encryptionPsk,omitempty"`
	GikRekeyPskEnable *bool   `json:"gikRekeyPskEnable,omitempty"`
	RekeyPskInterval  *int    `json:"rekeyPskInterval,omitempty"`
	IntervalPskType   *int    `json:"intervalPskType,omitempty"`
}

type WpaSetting struct {
	VersionEnt      *int    `json:"versionEnt,omitempty"`
	EncryptionEnt   *int    `json:"encryptionEnt,omitempty"`
	GikRekeyEnable  *bool   `json:"gikRekeyEnable,omitempty"`
	RekeyInterval   *bool   `json:"rekeyInterval,omitempty"`
	IntervalType    *int    `json:"intervalType,omitempty"`
	RadiusProfileId *string `json:"radiusProfileId,omitempty"`
}

type RateLimit struct {
	DownLimitEnable *bool   `json:"downLimitEnable,omitempty"`
	DownLimit       *int    `json:"downLimit,omitempty"`
	DownLimitType   *int    `json:"downLimitType,omitempty"`
	UpLimitEnable   *bool   `json:"upLimitEnable,omitempty"`
	UpLimit         *int    `json:"upLimit,omitempty"`
	UpLimitType     *int    `json:"upLimitType,omitempty"`
	RateLimitId     *string `json:"rateLimitId,omitempty"`
}

type RateAndBeaconCtrl struct {
	Rate2gCtrlEnable *bool `json:"rate2gCtrlEnable,omitempty"`
	Rate5gCtrlEnable *bool `json:"rate5gCtrlEnable,omitempty"`
	Rate6gCtrlEnable *bool `json:"rate6gCtrlEnable,omitempty"`
}

type OmadaTimeRangeData struct {
	Name     *string     `json:"name,omitempty"`
	DayMode  *int        `json:"dayMode,omitempty"`
	TimeList *[]TimeList `json:"timeList,omitempty"`
}

type NetworksSsidRequest struct {
	RadioOn       *bool        `json:"radioOn,omitempty"`
	UploadLimit   *int         `json:"uploadLimit,omitempty"`
	DownloadLimit *int         `json:"downloadLimit,omitempty"`
	UrlFilters    *[]UrlFilter `json:"urlFilters,omitempty"`
}

type NetworksResponse struct {
	Ssid          *string      `json:"ssid,omitempty"`
	RadioOn       *bool        `json:"radioOn,omitempty"`
	UploadLimit   *int         `json:"uploadLimit,omitempty"`
	DownloadLimit *int         `json:"downloadLimit,omitempty"`
	UrlFilters    *[]UrlFilter `json:"urlFilters,omitempty"`
	Updated       *bool        `json:"updated,omitempty"`
	ErrorMessage  *string      `json:"errorMessage,omitempty"`
}

type UrlFilter struct {
	Name   *string   `json:"name,omitempty"`
	Enable *bool     `json:"enable,omitempty"`
	Urls   *[]string `json:"urls,omitempty"`
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
