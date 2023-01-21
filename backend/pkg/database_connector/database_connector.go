package databaseConnector

import (
	dg "backend/pkg/data_gatherer"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DbConnector struct {
	getEndpoint       string
	preferredEndpoint string
	queryOptions      any
}

type PocketbaseGetQueryOptions struct {
	page    int    // !TODO not implemented
	perPage int    // !TODO not implemented
	sort    string // !TODO not implemented
	filter  string // !TODO not implemented
	expand  string // !TODO not implemented
}

type PocketbaseResponse struct {
	page       int
	perPage    int
	totalItems int
	totalPages int
	items      []dg.PowerPriceEntry
}

type GenericPowerEntry struct {
	HourUTC      string
	SpotPriceEUR float64
}

func (dbc DbConnector) _pocketbaseGetQuery(queryOptions any) ([]GenericPowerEntry, error) {

	// Executing an HTTP request
	resp, err := http.Get(dbc.getEndpoint)
	if err != nil {
		return []GenericPowerEntry{}, errors.New(fmt.Sprintf(`Unable to retrieve data from endpoint, got the following error: %v`, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []GenericPowerEntry{}, errors.New("Problems reading json response body from API call")
	}

	var jsonResponse map[string]any
	json.Unmarshal(body, &jsonResponse)
	if jsonResponse["code"] == float64(404) {
		return []GenericPowerEntry{}, errors.New("Error: got 404 from server (hint: improperly formed query or maybe database is down?")
	}

	return []GenericPowerEntry{
		{
			HourUTC:      "bla bla",
			SpotPriceEUR: 250.0,
		},
	}, nil

}

func (dbc DbConnector) GetQuery() ([]GenericPowerEntry, error) {
	switch dbc.preferredEndpoint {
	case "Pocketbase":
		switch dbc.queryOptions.(type) {
		case PocketbaseGetQueryOptions:
			res, err := dbc._pocketbaseGetQuery(dbc.queryOptions)
			if err != nil {
				return []GenericPowerEntry{}, errors.New(fmt.Sprintf(`Error executing get query on Pocketbase databse, got: %v`, err))
			}
			return res, nil
		}
	}
	return []GenericPowerEntry{}, errors.New("Unable to determine database, please visit: xxx to see valid database options") // !TODO documentation here
}
