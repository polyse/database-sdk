package database_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var healthcheckPath = "%s/healthcheck"
var contentType = "application/json"
var apiPath = "%s/api/%s/documents"
var DatabasePingErr = errors.New("can not ping database")

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

type ResponseData struct {
	Source
	Url string `json:"url"`
}

// DBClient struct consist payload and collection name
type DBClient struct {
	url string
	c   *http.Client
}

type CustomError struct {
	error
	code int
}

func wrap(msg string, code int, err error) error {
	return CustomError{error: fmt.Errorf("unexpected err: %w, body %s", err, msg), code: code}
}

// NewDBClient return new instance of DBClient
func NewDBClient(url string) (*DBClient, error) {
	db := &DBClient{
		url: url,
		c: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
	resp, err := db.c.Get(fmt.Sprintf(healthcheckPath, url))
	if err != nil {
		return nil, DatabasePingErr
	}
	if resp.StatusCode != http.StatusOK {
		return nil, DatabasePingErr
	}
	return db, nil
}

// SaveData RawData to PolySE Database
func (d *DBClient) SaveData(collectionName string, data Documents) (*Documents, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("can't perform request: %w", err)
	}
	resp, err := d.c.Post(
		fmt.Sprintf(
			apiPath,
			d.url,
			url.PathEscape(collectionName),
		),
		contentType,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
	}()
	var res Documents
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return &res, nil
	}
	if raw, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		return nil, wrap(string(raw), resp.StatusCode, err)
	}
}

// GetData returns data from PolySE Database
func (d *DBClient) GetData(collectionName, searchPhrase string, limit, offset int) ([]ResponseData, error) {
	response, err := d.c.Get(
		fmt.Sprintf(
			apiPath+"?q=%s&limit=%d&offset=%d",
			d.url,
			url.PathEscape(collectionName),
			url.PathEscape(searchPhrase),
			limit,
			offset,
		),
	)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, wrap(string(raw), response.StatusCode, err)
	}
	var result []ResponseData
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
