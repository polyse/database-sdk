package database_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Data string `json:"Data"`
}

// PolySEDB struct consist payload and collection name
type PolySEDB struct {
	Data           []RawData
	collectionName string
}

// New return new instance of PolySEDB
func New(d []RawData, name string) *PolySEDB {
	return &PolySEDB{
		Data:           d,
		collectionName: name,
	}
}

// Add RawData to PolySE Database
func (p *PolySEDB) Add() error {
	requestBody, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}
	resp, err := http.Post("/api/"+p.collectionName+"/documents", "application/json", bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		var answer string
		_ = json.NewDecoder(resp.Body).Decode(&answer)
		return errors.New(answer)
	}
}
