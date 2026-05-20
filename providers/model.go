package providers

import "github.com/google/uuid"

// Provider model that represents providers
type Provider struct {
	ID        uuid.UUID  `json:"_id"`
	Type      *string    `json:"type"`
	NJTransit *NJTransit `json:"njtransit,omitempty"`
}

const (
	// NJTransitType one of supported Provider types
	NJTransitType = "NJTransit"
)

// NJTransit model that represents NJ Transit provider
type NJTransit struct {
	OrigStationCode *string `json:"orig_station_code"`
	DestStationCode *string `json:"dest_station_code"`
}
