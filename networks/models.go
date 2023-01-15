package networks

// nested within sbserver response
type ControllerIdResponse struct {
	ErrorCode int                `json:"errorCode"`
	Msg       string             `json:"msg"`
	Result    ControllerIdResult `json:"result"`
}

type ControllerIdResult struct {
	OmadacId string `json:"omadacId"`
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
