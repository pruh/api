package models

import "time"

import "github.com/google/uuid"

// Timestamp represents unix time
type Timestamp time.Time

// Notification model for storing notifications
type Notification struct {
	ID        uuid.UUID
	Message   string
	StartTime Timestamp
	EndTime   Timestamp
	Source    string
}
