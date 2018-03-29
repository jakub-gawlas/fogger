package notifier

import (
	"time"
)

func WithInterval(interval time.Duration) func(*Notifier) {
	return func(notif *Notifier) {
		notif.Interval = interval
	}
}

func WithStore(store KVStore) func(*Notifier) {
	return func(notif *Notifier) {
		notif.Store = store
	}
}
