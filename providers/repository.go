package providers

import (
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

// Repository to access providers store
type Repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]Provider
}

// NewRepository creates a providers repository backed by in-memory storage.
func NewRepository() *Repository {
	return &Repository{
		data: map[uuid.UUID]Provider{},
	}
}

// GetAll returns all providers
func (r *Repository) GetAll() ([]Provider, error) {
	glog.Info("Querying for all providers")
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := []Provider{}
	for _, prov := range r.data {
		providers = append(providers, prov)
	}
	return providers, nil
}

// GetOne returns provider by ID or nil
func (r *Repository) GetOne(providerID uuid.UUID) (*Provider, error) {
	glog.Infof("Querying for provider with UUID: %s", providerID)
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.data[providerID]
	if !ok {
		return nil, nil
	}

	return &provider, nil
}

// CreateOne creates new provider for specified params
func (r *Repository) CreateOne(provider Provider) bool {
	glog.Infof("Creating new provider: %+v", provider)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[provider.ID] = provider
	return true
}

// DeleteOne deletes provider with ID and returns true if record was removed
func (r *Repository) DeleteOne(providerID uuid.UUID) (bool, error) {
	glog.Infof("Deleting provider with UUID: %s", providerID)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[providerID]; !ok {
		return false, nil
	}
	delete(r.data, providerID)

	return true, nil
}
