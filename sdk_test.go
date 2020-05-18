package database_sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	server *httptest.Server
	client *DBClient
	res    []ResponseData
	req    []RawData
}

func TestSetupServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	t := time.Now().Round(1 * time.Nanosecond)
	s.res = []ResponseData{
		{
			Source: Source{
				Date:  t,
				Title: "Test title",
			},
			Url: "http://testurl.com",
		},
	}
	s.req = []RawData{{
		Source: Source{
			Date:  t,
			Title: "Test title",
		},
		Url:  "http://testurl.com",
		Data: "a b c d",
	}}
	b, err := json.Marshal(s.res)
	if err != nil {
		panic(err)
	}

	r := http.NewServeMux()
	r.HandleFunc("/api/default/documents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if _, err := fmt.Fprint(w, string(b)); err != nil {
				panic(err)
			}
			return
		}
		if r.Method == http.MethodPost {
			raw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			if _, err := fmt.Fprint(w, string(raw)); err != nil {
				panic(err)
			}
			return
		}
		panic("method not allowed")

	})
	s.server = httptest.NewServer(r)
	s.client = NewDBClient(s.server.URL)
}

func (s *ServerTestSuite) TestClient_Get() {
	d, err := s.client.GetData("default", "data1", 10, 0)
	s.NoError(err)
	s.ElementsMatch(s.res, d)
}

func (s *ServerTestSuite) TestClient_Post() {
	d, err := s.client.SaveData("default", Documents{Documents: s.req})
	s.NoError(err)
	s.ElementsMatch(s.req, d.Documents)
}

func (s *ServerTestSuite) TearDownSuite() {
	s.server.Close()
}
