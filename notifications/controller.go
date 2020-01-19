package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pruh/api/mongo"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	uuidParam   = "uuid"
	uuidPattern = "{" + uuidParam + "}"
	// GetPath path for HTTP GET Method which to query for array of notifications
	GetPath = "/notifications"
	// SingleGetPath path for HTTP GET Method to query for single notification
	SingleGetPath = GetPath + "/" + uuidPattern
	// CreatePath path for HTTP POST Method which creates new notification
	CreatePath = "/notifications"
	// DeletePath path for HTTP DELETE Method which deletes notification
	DeletePath = "/notifications/" + uuidPattern
)

// Controller handles all notification related requests.
type Controller struct {
	Repository *Repository
}

// GetAll returns all notifications.
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	notifications, err := c.Repository.GetAll()
	if err != nil {
		glog.Errorf("Error while querying notifications. %s", err)
		http.Error(w, fmt.Sprint("Error while querying notifications."), http.StatusInternalServerError)
		return
	}

	current, ok := r.URL.Query()["only_current"]
	if ok && "true" == current[0] {
		notifications = FilterNotificatons(notifications, CurrentFilter)
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
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notifUUID := params[uuidParam]
	mongoUUID, err := validateUUID(notifUUID)
	if err != nil {
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	notification, err := c.Repository.GetOne(*mongoUUID)
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
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var notification Notification
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

	notification.ID = mongo.NewUUID()
	success := c.Repository.CreateOne(notification)
	if !success {
		glog.Errorf("Failed to create notification. %s", err)
		http.Error(w, fmt.Sprint("Failed to create notification."), http.StatusInternalServerError)
		return
	}

	// add relative location header as required by RFC 7231 § 7.1.2
	getPath := strings.ReplaceAll(SingleGetPath, uuidPattern, notification.ID.String())
	location := strings.ReplaceAll(r.URL.String(), CreatePath, getPath)
	w.Header().Set("Location", location)

	w.WriteHeader(http.StatusCreated)
	return
}

// Delete deletes notification.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	notifUUID := params[uuidParam]
	mongoUUID, err := validateUUID(notifUUID)
	if err != nil {
		glog.Errorf("Notification UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Notification UUID is malformed."), http.StatusBadRequest)
		return
	}

	res, err := c.Repository.DeleteOne(*mongoUUID)
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

func validateUUID(notifUUID string) (*mongo.UUID, error) {
	u, err := uuid.Parse(notifUUID)
	if err != nil {
		return nil, err
	}

	return &mongo.UUID{UUID: u}, nil
}

func validateNotification(notif Notification) error {
	if notif.Title == nil {
		return errors.New("title should be set")
	}
	if notif.StartTime.IsZero() {
		return errors.New("start_time should be set")
	}
	if notif.EndTime.IsZero() {
		return errors.New("end_time should be set")
	}
	if notif.StartTime.After(notif.EndTime.Time) {
		return errors.New("start_time can not be after end_time")
	}
	if notif.EndTime.Before(time.Now()) {
		return errors.New("end_time is in the past")
	}
	return nil
}
