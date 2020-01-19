package notifications

import "github.com/pruh/api/mongo"

// Notification model for storing notifications
type Notification struct {
	ID        mongo.UUID       `json:"_id" bson:"_id"`
	Title     *string          `json:"title" bson:"title"`
	Message   *string          `json:"message,omitempty" bson:"message,omitempty"`
	StartTime *mongo.Timestamp `json:"start_time" bson:"start_time"`
	EndTime   *mongo.Timestamp `json:"end_time" bson:"end_time"`
	Source    *string          `json:"source,omitempty" bson:"source,omitempty"`
}
