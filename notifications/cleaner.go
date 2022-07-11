package notifications

import (
	"time"

	"github.com/golang/glog"
	"github.com/pruh/api/mongo"
)

const (
	// period how often to fire cleaner
	period = time.Minute
)

// Cleaner handles all work related to cleaning repository.
type Cleaner struct {
	Repository *Repository
}

// StartPeriodicCleaner starts periodic task to cleanup
func (c *Cleaner) StartPeriodicCleaner() {
	ticker := time.NewTicker(period)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				c.removeExpired()
			}
		}
	}()
}

func (c *Cleaner) removeExpired() {
	glog.Info("Checking for expired notifications")
	notifs, err := c.Repository.GetAll()
	if err != nil {
		glog.Error("Cannot query for notifications. ", err)
		return
	}

	if len(notifs) == 0 {
		glog.Info("No notifications found")
		return
	}

	notifs = FilterNotificatons(notifs, ExpiredFilter)
	if len(notifs) == 0 {
		glog.Info("No expired notifications found")
		return
	}

	glog.Infof("Found expired notifications %+v", notifs)

	ids := extractUuids(notifs)

	glog.Info("Deleting expired notifications with ids: ", ids)
	_, err = c.Repository.DeleteAll(ids)
	if err != nil {
		glog.Error("Cannot delete expired notifications. ", err)
		return
	}
}

func extractUuids(notifs []Notification) (ids []mongo.UUID) {
	for _, notif := range notifs {
		ids = append(ids, notif.ID)
	}
	return ids
}
