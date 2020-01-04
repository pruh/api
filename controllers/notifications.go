package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pruh/api/dao"
	"github.com/pruh/api/models"
)

// NotificationsController handles all notification related requests.
type NotificationsController struct {
	Repository *dao.Repository
}

// GetAll returns all notifications.
func (c *NotificationsController) GetAll(w http.ResponseWriter, r *http.Request) {
	notifications, err := c.Repository.GetNofitications()
	if err != nil {
		glog.Errorf("Error while querying notifications. %s", err)
		http.Error(w, fmt.Sprint("Error while querying notifications."), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(notifications)
	if err != nil {
		glog.Errorf("Cannot marshal notifications. %s", err)
		http.Error(w, fmt.Sprint("Cannot marshal notifications."), http.StatusInternalServerError)
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
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	notification := c.Repository.GetNofitication(notifUUID)
	data, err := json.Marshal(notification)
	if err != nil {
		glog.Errorf("Cannot marshal notification. %s", err)
		http.Error(w, fmt.Sprint("Cannot marshal notification."), http.StatusInternalServerError)
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
		glog.Errorf("Error reading request. %s", err)
		http.Error(w, fmt.Sprint("Error reading request."), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &notification)
	if err != nil {
		glog.Errorf("Error reading request data. %s", err)
		http.Error(w, fmt.Sprint("Error reading request data."), http.StatusUnprocessableEntity)
		return
	}

	success := c.Repository.CreateNofitication(notification)
	if !success {
		glog.Errorf("Failed to create notification. %s", err)
		http.Error(w, fmt.Sprint("Failed to create notification."), http.StatusInternalServerError)
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
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	res := c.Repository.DeleteNofitication(notifUUID)
	if !res {
		glog.Infof("Notification with ID %s was not removed", notifUUID)
		http.Error(w, fmt.Sprintf("Notification with ID %s was not removed", notifUUID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func validateUUID(notifUUID string) (uuid.UUID, error) {
	return uuid.Parse(notifUUID)
}
