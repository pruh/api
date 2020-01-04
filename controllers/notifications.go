package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pruh/api/dao"
	"github.com/pruh/api/models"
)

// NotificationsController handles all notification related requests.
type NotificationsController struct {
	Repository dao.Repository
}

// GetAll returns all notifications.
func (c *NotificationsController) GetAll(w http.ResponseWriter, r *http.Request) {
	notifications := c.Repository.GetNofitications()
	data, err := json.Marshal(notifications)
	if err != nil {
		log.Fatalln("Cannot marshal notifications.", err)
		http.Error(w, fmt.Sprintf("Cannot marshal notifications: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
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

	notification := c.Repository.GetNofitication(notifUUID)
	data, err := json.Marshal(notification)
	if err != nil {
		log.Fatalln("Cannot marshal notification.", err)
		http.Error(w, fmt.Sprintf("Cannot marshal notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// Create creates a new notification.
func (c *NotificationsController) Create(w http.ResponseWriter, r *http.Request) {
	var notification models.Notification
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		log.Fatalln("Error reading request.", err)
		http.Error(w, fmt.Sprintf("Error reading request: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &notification)
	if err != nil {
		log.Fatalln("Error reading request data.", err)
		http.Error(w, fmt.Sprintf("Error reading request data: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	success := c.Repository.CreateNofitication(notification)
	if !success {
		log.Fatalln("Failed to create notification.", err)
		http.Error(w, fmt.Sprintf("Failed to create notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
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

	res := c.Repository.DeleteNofitication(notifUUID)
	if !res {
		log.Printf("Notification with ID %s was not removed", notifUUID)
		http.Error(w, fmt.Sprintf("Notification with ID %s was not removed", notifUUID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func validateUUID(notifUUID string) (uuid.UUID, error) {
	return uuid.Parse(notifUUID)
}
