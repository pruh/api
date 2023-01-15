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

type LoginData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
