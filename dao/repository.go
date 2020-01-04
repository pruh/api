package dao

import (
	"log"

	"github.com/pruh/api/models"
)

// Repository accesses notifications store
type Repository struct{}

// GetNofitications returns all notifications
func (r *Repository) GetNofitications() []models.Notification {
	log.Printf("Querying for all notifications")
	return []models.Notification{}
}

// GetNofitication returns notifications by ID or nil
func (r *Repository) GetNofitication(ID string) models.Notification {
	log.Printf("Querying for notification with UUID: %s\n", ID)
	return models.Notification{}
}

// CreateNofitication creates new notification for specified params
func (r *Repository) CreateNofitication(notification models.Notification) bool {
	log.Printf("Creating new notification: %+v\n", notification)
	return true
}

// DeleteNofitication deletes notifications with ID
func (r *Repository) DeleteNofitication(ID string) bool {
	log.Printf("Deleting notification with UUID: %s\n", ID)
	// todo handle not found vs internal error
	return true
}
