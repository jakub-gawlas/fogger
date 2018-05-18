package client

import (
	"bytes"
	"encoding/json"
	"net/url"
)

// Entry represents journal entry
type Entry struct {
	Hash   string `json:"hash,omitempty"`
	Links  Links  `json:"links,omitempty"`
	Object string `json:"object,omitempty"`
}

// Links represents links for entries
type Links struct {
	PreviousEntry string `json:"previousEntry,omitempty"`
}

// GetJournalTail returns last entry from journal
func (c *Client) GetJournalTail(journal string) (*Entry, error) {
	q, _ := url.Parse("journals/" + journal + "/tail")
	path := c.URL.ResolveReference(q).String()

	r, err := c.http.Get(path)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	entry := Entry{}

	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		return nil, nil
	}

	return &entry, nil
}

// GetJournalEntry returns entry with the given hash
func (c *Client) GetJournalEntry(hash string) (*Entry, error) {
	q, _ := url.Parse("journals/" + hash + "/json")
	path := c.URL.ResolveReference(q).String()

	r, err := c.http.Get(path)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	entry := Entry{}

	err = json.NewDecoder(r.Body).Decode(&entry)

	return &entry, nil
}

// GetEntries returns limited number (if limit != 0) of entries until given hash
func (c *Client) GetEntries(journal, endHash string, limit int) ([]*Entry, error) {
	entries := []*Entry{}
	hash := ""
	i := 0

	for {
		if limit != 0 && i == limit {
			return entries, nil
		}

		var (
			entry *Entry
			err   error
		)

		if hash == "" {
			entry, err = c.GetJournalTail(journal)
		} else {
			entry, err = c.GetJournalEntry(hash)
		}

		if err != nil {
			return nil, err
		}

		if entry.Hash == endHash {
			return entries, nil
		}

		entries = append(entries, entry)

		if entry.Links.PreviousEntry == "" {
			return entries, nil
		}

		hash = entry.Links.PreviousEntry
		i++
	}
}

// PushToJournal adds to journal file with the given hash
func (c *Client) PushToJournal(journal, hash string) error {
	q, _ := url.Parse("journals/" + journal + "/push")
	path := c.URL.ResolveReference(q).String()

	body := map[string]string{"object": hash}
	jsonBody, _ := json.Marshal(body)

	r, err := c.http.Post(path, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

// CreateJournal creates new journal
func (c *Client) CreateJournal(journal string) error {
	q, _ := url.Parse("journals")
	path := c.URL.ResolveReference(q).String()

	body := map[string]string{"name": journal}
	jsonBody, _ := json.Marshal(body)

	r, err := c.http.Post(path, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}
