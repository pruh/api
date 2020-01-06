package models

// Notification model for storing notifications
type Notification struct {
	ID        MongoUUID  `json:"_id" bson:"_id"`
	Title     *string    `json:"title" bson:"title"`
	Message   *string    `json:"message,omitempty" bson:"message,omitempty"`
	StartTime *Timestamp `json:"start_time" bson:"start_time"`
	EndTime   *Timestamp `json:"end_time" bson:"end_time"`
	Source    *string    `json:"source,omitempty" bson:"source,omitempty"`
}
