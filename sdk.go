package database_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Временно
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
}

// NewDBClient return new instance of DBClient
func NewDBClient(url string) *DBClient {
	return &DBClient{url: url}
}

// SaveData RawData to PolySE Database
func (d *DBClient) SaveData(collectionName string, data Documents) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/api/%s/documents", d.url, collectionName), contentType, bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusOK {
		return nil
	}
	var answer string
	_ = json.NewDecoder(resp.Body).Decode(&answer)
	return errors.New(answer)
}
