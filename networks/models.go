package networks

// nested within sbserver response
type OmadaResponse struct {
	ErrorCode int    `json:"errorCode"`
	Msg       string `json:"msg"`
	Result    Result `json:"result"`
}

type Result struct {
	OmadacId *string `json:"omadacId,omitempty"`
	Token    *string `json:"token,omitempty"`
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
