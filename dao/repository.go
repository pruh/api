package dao

import (
	"github.com/golang/glog"

	"github.com/pruh/api/models"
)

// Repository accesses notifications store
type Repository struct{}

// GetNofitications returns all notifications
func (r *Repository) GetNofitications() []models.Notification {
	glog.Infof("Querying for all notifications")
	return []models.Notification{}
}

// GetNofitication returns notifications by ID or nil
func (r *Repository) GetNofitication(ID string) models.Notification {
	glog.Infof("Querying for notification with UUID: %s\n", ID)
	return models.Notification{}
}

// CreateNofitication creates new notification for specified params
func (r *Repository) CreateNofitication(notification models.Notification) bool {
	glog.Infof("Creating new notification: %+v\n", notification)
	return true
}

// DeleteNofitication deletes notifications with ID
func (r *Repository) DeleteNofitication(ID string) bool {
	glog.Infof("Deleting notification with UUID: %s\n", ID)
	// todo handle not found vs internal error
	return true
}
