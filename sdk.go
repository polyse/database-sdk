package database_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Source struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
}
type RawData struct {
	Source
	Url  string `json:"url"`
	Data string `json:"data"`
}

type Documents struct {
	Documents []RawData `json:"documents"`
}

var contentType = "application/json"

// DBClient struct consist payload and collection name
type DBClient struct {
	url string
	c   *http.Client
}

// NewDBClient return new instance of DBClient
func NewDBClient(url string) *DBClient {
	return &DBClient{
		url: url,
		c: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// SaveData RawData to PolySE Database
func (d *DBClient) SaveData(collectionName string, data Documents) (*Documents, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Can't perform request: %w", err)
	}
	resp, err := d.c.Post(fmt.Sprintf("%s/api/%s/documents", d.url, collectionName), contentType, bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res := struct {
		D       Documents `json:"documents"`
		Message string    `json:"message"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("Unexpected answer: %w", err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return &res.D, nil
	}
	if res.Message != "" {
		return &res.D, errors.New(res.Message)
	}
	return nil, fmt.Errorf("Unexpected answer: %w", err)
}
