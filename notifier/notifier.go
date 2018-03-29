package notifier

import (
	"time"

	"github.com/jakub-gawlas/fogger/client"

	log "github.com/sirupsen/logrus"
)

type Notifier struct {
	Store     KVStore
	Interval  time.Duration
	listeners []*Listener
}

type KVStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

func New(modifiers ...func(*Notifier)) *Notifier {
	notif := &Notifier{
		Interval: time.Minute,
		Store:    newStore(),
	}
	for _, modifier := range modifiers {
		modifier(notif)
	}
	return notif
}

// Add new listener
func (n *Notifier) Add(listener *Listener) *Notifier {
	n.listeners = append(n.listeners, listener)

	return n
}

// Start notifier
func (n *Notifier) Start() {
	n.tick(1)
	ticker := time.NewTicker(n.Interval)
	for _ = range ticker.C {
		n.tick(0)
	}
}

func (n *Notifier) tick(limit int) {
	for _, l := range n.listeners {
		id := l.ID()
		hash, err := n.Store.Get(id)
		if err != nil {
			log.Errorf("while get hash: %v", err)
		}

		fog := client.New(l.URL)
		entries, err := fog.GetEntries(l.Journal, hash, limit)
		if err != nil {
			log.Errorf("while get entries: %v", err)
		}

		if len(entries) == 0 {
			continue
		}

		if l.Events.OnEntryAdded != nil {
			for i := range entries {
				go l.Events.OnEntryAdded(entries[len(entries)-i-1])
			}
		}

		last := entries[len(entries)-1]
		if err := n.Store.Set(id, last.Hash); err != nil {
			log.Errorf("while set hash: %v", err)
		}
	}
}
