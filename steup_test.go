package main

import (
	"bytes"
	"database/sql"
	"fyne.io/fyne/v2/test"
	"gold/repository"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

var testApp Config

func TestMain(m *testing.M) {
	a := test.NewApp()
	testApp.App = a
	testApp.HTTPClient = client
	_ = os.Remove("./repository/testdata/sql.db")
	path := "./repository/testdata/sql.db"
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Println(err)
	}
	testApp.DB = repository.NewSQLiteRepository(db)
	err = testApp.DB.Migrate()
	if err != nil {
		log.Println("migrate failed")
	}
	os.Exit(m.Run())
}

var jsonToReturn = `
{
    "ts": 1667991959681, 
    "tsj": 1667991955958, 
    "date": "Nov 9th 2022, 06:05:55 am NY", 
    "items": [
        {
            "curr": "USD", 
            "xauPrice": 1709.2, 
            "xagPrice": 21.3665, 
            "chgXau": -2.405, 
            "chgXag": -0.0175, 
            "pcXau": -0.1405, 
            "pcXag": -0.0818, 
            "xauClose": 1711.605, 
            "xagClose": 21.384
        }
    ]
}
`

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

var client = NewTestClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
		Header:     make(http.Header),
	}
})
