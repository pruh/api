package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"time"

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
	mongoUUID, err := validateUUID(notifUUID)
	if err != nil {
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	notification, err := c.Repository.GetNofitication(*mongoUUID)
	if err != nil {
		glog.Errorf("Error while querying notification. %s", err)
		http.Error(w, fmt.Sprint("Error while querying notification."), http.StatusInternalServerError)
		return
	}
	if notification == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
	notification := models.NewNotification()
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

	err = validateNotification(notification)
	if err != nil {
		glog.Error(err)
		http.Error(w, fmt.Sprintf("Input data is not valid. %s", err), http.StatusUnprocessableEntity)
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
	mongoUUID, err := validateUUID(notifUUID)
	if err != nil {
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	res, err := c.Repository.DeleteNofitication(*mongoUUID)
	if err != nil {
		glog.Errorf("Failed to delete notification. %s", err)
		http.Error(w, fmt.Sprintf("Failed to delete notification. %s", err), http.StatusInternalServerError)
		return
	}

	if !res {
		glog.Infof("Notification with ID %s was not removed", notifUUID)
		http.Error(w, fmt.Sprintf("Notification with ID %s was not removed", notifUUID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func validateUUID(notifUUID string) (*models.MongoUUID, error) {
	u, err := uuid.Parse(notifUUID)
	if err != nil {
		return nil, err
	}

	return &models.MongoUUID{UUID: u}, nil
}

func validateNotification(notif models.Notification) error {
	if len(notif.Message) == 0 {
		return errors.New("message not set")
	}
	if notif.StartTime.IsZero() {
		return errors.New("start_time not set")
	}
	if notif.EndTime.IsZero() {
		return errors.New("end_time not set")
	}
	if notif.StartTime.After(notif.EndTime.Time) {
		return errors.New("start_time can not be after end_time")
	}
	if notif.EndTime.Before(time.Now()) {
		return errors.New("end_time is in the past")
	}
	return nil
}
