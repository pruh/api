package notifications

import (
	"time"

	"github.com/pruh/api/notifications/models"
)

// FilterNotificatons filters notification using provided filter func
func FilterNotificatons(notifications []models.Notification,
	filterFunc func(models.Notification) bool) []models.Notification {
	filtered := []models.Notification{}
	for _, notif := range notifications {
		if filterFunc(notif) {
			filtered = append(filtered, notif)
		}
	}
	return filtered
}

// CurrentFilter returns true if notification is current
func CurrentFilter(notif models.Notification) bool {
	now := time.Now()
	return notif.StartTime.Before(now) && now.Before(notif.EndTime.Time)
}

// ExpiredFilter returns true if notification is expired
func ExpiredFilter(notif models.Notification) bool {
	return notif.EndTime.Before(time.Now())
}
