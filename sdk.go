package database_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
func (d *DBClient) SaveData(collectionName string, data Documents) (Documents, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return data, err
	}
	client := http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Post(fmt.Sprintf("%s/api/%s/documents", d.url, collectionName), contentType, bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	if err != nil {
		return data, err
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return data, errors.New(http.StatusText(resp.StatusCode))
	}
	res := struct {
		D       Documents `json:"documents"`
		Message string    `json:"message"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return data, err
	}
	if res.Message != "" {
		if strings.Contains(res.Message, "200") == true || strings.Contains(res.Message, "201") == true {
			return res.D, nil
		}
		return res.D, errors.New(res.Message)
	}
	return res.D, errors.New("Unexpected answer")
}
