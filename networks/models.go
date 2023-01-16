package networks

// nested within sbserver response
type OmadaResponse struct {
	ErrorCode int     `json:"errorCode"`
	Msg       *string `json:"msg,omitempty"`
	Result    *Result `json:"result,omitempty"`
}

type Result struct {
	OmadacId *string `json:"omadacId,omitempty"`
	Token    *string `json:"token,omitempty"`
	Data     *[]Data `json:"data,omitempty"`
}

type Data struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
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
