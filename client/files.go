package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/jakub-gawlas/go-retryablehttp"
)

// GetFile returns file with the given hash
func (c *Client) GetFile(hash string) ([]byte, error) {
	q, _ := url.Parse("files/" + hash)
	path := c.URL.ResolveReference(q).String()

	req, err := retryablehttp.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/octet-stream")

	r, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetTailFile returns file for the tail
func (c *Client) GetTailFile(journal string) ([]byte, error) {
	entry, err := c.GetJournalTail(journal)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("no tail")
	}

	if entry.Object == "" {
		return nil, nil
	}

	return c.GetFile(entry.Object)
}

// AddFileRes is response for upload file
type AddFileRes struct {
	Hash string `json:"hash"`
}

// AddFile uploads file to fogger
func (c *Client) AddFile(file io.ReadSeeker) (*AddFileRes, error) {
	q, _ := url.Parse("files")
	path := c.URL.ResolveReference(q).String()

	req, err := retryablehttp.NewRequest("POST", path, file)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	r, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	dto := AddFileRes{}
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

// AddFileAndPush uploads file and push the one to journal
func (c *Client) AddFileAndPush(journal string, file io.ReadSeeker) error {
	fileStat, err := c.AddFile(file)
	if err != nil {
		return err
	}

	err = c.PushToJournal(journal, fileStat.Hash)
	return err
}
