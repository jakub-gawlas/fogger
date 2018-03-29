package notifier

import (
	"fmt"
	"net/url"

	"github.com/jakub-gawlas/fogger/client"
)

type Listener struct {
	URL     *url.URL
	Journal string
	Events  *Events
}

type Events struct {
	OnEntryAdded Handler
}

type Handler func(*client.Entry)

func NewListener(rawurl, journal string, onEvents ...func(*Events)) (*Listener, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("invalid fogger url: %v", err)
	}

	e := &Events{}
	for _, on := range onEvents {
		on(e)
	}

	l := &Listener{
		URL:     u,
		Journal: journal,
		Events:  e,
	}

	return l, nil
}

func OnEntryAdded(handler Handler) func(*Events) {
	return func(e *Events) {
		e.OnEntryAdded = handler
	}
}

func (l *Listener) ID() string {
	return fmt.Sprintf("%s-%s", l.URL.Path, l.Journal)
}
