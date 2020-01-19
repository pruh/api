package notifications

import (
	"time"
)

// FilterNotificatons filters notification using provided filter func
func FilterNotificatons(notifications []Notification,
	filterFunc func(Notification) bool) []Notification {
	filtered := []Notification{}
	for _, notif := range notifications {
		if filterFunc(notif) {
			filtered = append(filtered, notif)
		}
	}
	return filtered
}

// CurrentFilter returns true if notification is current
func CurrentFilter(notif Notification) bool {
	now := time.Now()
	return notif.StartTime.Before(now) && now.Before(notif.EndTime.Time)
}

// ExpiredFilter returns true if notification is expired
func ExpiredFilter(notif Notification) bool {
	return notif.EndTime.Before(time.Now())
}
