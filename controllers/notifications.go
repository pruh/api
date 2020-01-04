package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// NotificationsController handles all notification related requests.
type NotificationsController struct {
}

// GetAll returns all notifications.
func (c *NotificationsController) GetAll(w http.ResponseWriter, r *http.Request) {

}

// Get returns one notification.
func (c *NotificationsController) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notifUUID := params["uuid"]
	_, err := validateUUID(notifUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Notification UUID is malformed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Printf("Querying for notifications with UUID: %s\n", notifUUID)
}

// Create creates a new notification.
func (c *NotificationsController) Create(w http.ResponseWriter, r *http.Request) {
}

// Delete deletes notification.
func (c *NotificationsController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notifUUID := params["uuid"]
	_, err := validateUUID(notifUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Notification UUID is malformed: %s", err.Error()), http.StatusBadRequest)
		return
	}

	log.Printf("Querying for notifications with UUID: %s\n", notifUUID)
}

func validateUUID(notifUUID string) ([16]byte, error) {
	return uuid.Parse(notifUUID)
}
