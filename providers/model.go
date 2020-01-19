package providers

import "github.com/pruh/api/mongo"

// Provider model that represents providers
type Provider struct {
	ID   mongo.UUID `json:"_id" bson:"_id"`
	Type *string    `json:"type" bson:"type"`

	// Title     *string    `json:"title" bson:"title"`
	// Message   *string    `json:"message,omitempty" bson:"message,omitempty"`
	// StartTime *Timestamp `json:"start_time" bson:"start_time"`
	// EndTime   *Timestamp `json:"end_time" bson:"end_time"`
	// Source    *string    `json:"source,omitempty" bson:"source,omitempty"`
}
