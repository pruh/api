package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pruh/api/mongo"
)

const (
	uuidParam   = "uuid"
	uuidPattern = "{" + uuidParam + "}"
	// GetPath path for HTTP GET Method which to query for providers
	GetPath = "/providers"
	// SingleGetPath path for HTTP GET Method to query for a single provider
	SingleGetPath = GetPath + "/" + uuidPattern
	// CreatePath path for HTTP POST Method to create a new provider
	CreatePath = "/providers"
	// DeletePath path for HTTP DELETE Method to delete a provider
	DeletePath = "/providers/" + uuidPattern
)

// Controller handles all provider related requests.
type Controller struct {
	Repository *Repository
}

// GetAll returns all providers.
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	providers, err := c.Repository.GetAll()
	if err != nil {
		glog.Errorf("Error while querying providers. %s", err)
		http.Error(w, fmt.Sprint("Error while querying providers."), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(providers)
	if err != nil {
		glog.Errorf("Cannot marshal providers. %s", err)
		http.Error(w, fmt.Sprint("Cannot marshal providers."), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// Get returns one provider.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuid := params[uuidParam]
	mongoUUID, err := validateUUID(uuid)
	if err != nil {
		glog.Errorf("Provider UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Provider UUID is malformed."), http.StatusBadRequest)
		return
	}

	provider, err := c.Repository.GetOne(*mongoUUID)
	if err != nil {
		glog.Errorf("Error while querying provider. %s", err)
		http.Error(w, fmt.Sprint("Error while querying provider."), http.StatusInternalServerError)
		return
	}
	if provider == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(provider)
	if err != nil {
		glog.Errorf("Cannot marshal provider. %s", err)
		http.Error(w, fmt.Sprint("Cannot marshal provider."), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// Create creates a new provider.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var provider Provider
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		glog.Errorf("Error reading request. %s", err)
		http.Error(w, fmt.Sprint("Error reading request."), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &provider)
	if err != nil {
		glog.Errorf("Error reading request data. %s", err)
		http.Error(w, fmt.Sprint("Error reading request data."), http.StatusUnprocessableEntity)
		return
	}

	err = validateProvider(provider)
	if err != nil {
		glog.Error(err)
		http.Error(w, fmt.Sprintf("Input data is not valid. %s", err), http.StatusUnprocessableEntity)
		return
	}

	provider.ID = mongo.NewUUID()
	success := c.Repository.CreateOne(provider)
	if !success {
		glog.Errorf("Failed to create provider. %s", err)
		http.Error(w, fmt.Sprint("Failed to create provider."), http.StatusInternalServerError)
		return
	}

	// add relative location header as required by RFC 7231 § 7.1.2
	getPath := strings.ReplaceAll(SingleGetPath, uuidPattern, provider.ID.String())
	location := strings.ReplaceAll(r.URL.String(), CreatePath, getPath)
	w.Header().Set("Location", location)

	w.WriteHeader(http.StatusCreated)
	return
}

// Delete deletes provider.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuid := params[uuidParam]
	mongoUUID, err := validateUUID(uuid)
	if err != nil {
		glog.Errorf("Provider UUID is malformed. %s", err)
		http.Error(w, fmt.Sprint("Provider UUID is malformed."), http.StatusBadRequest)
		return
	}

	res, err := c.Repository.DeleteOne(*mongoUUID)
	if err != nil {
		glog.Errorf("Failed to delete provider. %s", err)
		http.Error(w, fmt.Sprintf("Failed to delete provider. %s", err), http.StatusInternalServerError)
		return
	}

	if !res {
		glog.Infof("Provider with ID %s was not removed", uuid)
		http.Error(w, fmt.Sprintf("Provider with ID %s was not removed", uuid), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func validateUUID(uuidStr string) (*mongo.UUID, error) {
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}

	return &mongo.UUID{UUID: u}, nil
}

func validateProvider(prov Provider) error {
	return nil
}
