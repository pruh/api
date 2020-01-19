package providers

import "github.com/pruh/api/mongo"

// Provider model that represents providers
type Provider struct {
	ID        mongo.UUID `json:"_id" bson:"_id"`
	Type      *string    `json:"type" bson:"type"`
	NJTransit *NJTransit `json:"njtransit,omitempty" bson:"njtransit,omitempty"`
}

const (
	// NJTransitType one of supported Provider types
	NJTransitType = "NJTransit"
)

// NJTransit model that represents NJ Transit provider
type NJTransit struct {
	StationID *string `json:"station_id" bson:"station_id"`
}
