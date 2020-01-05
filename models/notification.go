package models

// Notification model for storing notifications
type Notification struct {
	ID        MongoUUID `json:"_id" bson:"_id"`
	Message   string    `json:"message" bson:"message"`
	StartTime Timestamp `json:"start_time" bson:"start_time"`
	EndTime   Timestamp `json:"end_time" bson:"end_time"`
	Source    *string   `json:"source,omitempty" bson:"source,omitempty"`
}
