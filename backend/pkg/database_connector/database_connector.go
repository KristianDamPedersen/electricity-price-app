package databaseConnector

import (
	dg "backend/pkg/data_gatherer"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"net/http"
)

// DbConnector is the main struct where we keep track of the database connection.
//
// endpoint is the endpoint to where get requests will be executed against.
//
// database is the desired database where contacting (only valid option at the moment is "Pocketbase")
//
// queryOptions expects a struct containing some parameters to be inserted into the query. Exact implementation will depend
// on which database we use. But the queryOptions will always be validated before any query is executed.
type DbConnector struct {
	queryType    string
	database     string
	queryOptions any
}

type PocketbaseGetQueryOptions struct {
	page    int    // !TODO not implemented
	perPage int    // !TODO not implemented
	sort    string // !TODO not implemented
	filter  string // !TODO not implemented
	expand  string // !TODO not implemented
}

type _pocketbaseResponse struct {
	page       int
	perPage    int
	totalItems int
	totalPages int
	items      []dg.EnergiDataServiceEntry
}

type GenericPowerEntry struct {
	datetimeUTC string
	priceEUR    float64
	priceDKK    float64
}

// HttpGetRequest executes a get request against a specific endpoint / query string.
//
// The queryString parameter is a full string containing both the address and associated query (so full HTTP req).
//
// Returns a []byte of the response, which can be unmarshalled later.
func HttpGetRequest(queryString string) ([]byte, error) {
	resp, err := http.Get(queryString)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf(`Unable to retrieve data from endpoint, got the following error: %v`, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Problems reading json response body from API call")
	}
	return body, nil
}

// Unmashals the json response from pocketbase into the _pocketbaseResponse type.
func _unmarshalJsonToPocketbaseResponse(jsonBody []byte) (_pocketbaseResponse, error) {
	var jsonResponse map[string]any
	json.Unmarshal(jsonBody, &jsonResponse)
	if jsonResponse["code"] == float64(404) {
		return _pocketbaseResponse{}, errors.New("Error: got 404 from server (hint: improperly formed query or maybe database is down?")
	}

	pocketResponse := _pocketbaseResponse{
		page:       int(jsonResponse["page"].(float64)),
		perPage:    int(jsonResponse["perPage"].(float64)),
		totalItems: int(jsonResponse["totalItems"].(float64)),
		totalPages: int(jsonResponse["totalPages"].(float64)),
		items:      []dg.EnergiDataServiceEntry{},
	}

	for i := range jsonResponse["items"].([]interface{}) {
		items := jsonResponse["items"].([]interface{})[i].(map[string]any)
		entry := dg.EnergiDataServiceEntry{
			HourUTC:      items["datetime_UTC"].(string),
			HourDK:       items["datetime_DK"].(string),
			PriceArea:    items["priceArea"].(string),
			SpotPriceDKK: items["priceDKK"].(float64),
			SpotPriceEUR: items["priceEUR"].(float64),
		}
		pocketResponse.items = append(pocketResponse.items, entry)
	}
	return pocketResponse, nil
}

// Specifically executes a get request against pocketbase.
func (dbc DbConnector) _pocketbaseGetQuery(queryOptions any) ([]GenericPowerEntry, error) {
	// Read TOML
	t, err := toml.LoadFile("../../config.toml")
	endpoint := t.Get("database.pocketbase.getEndpoint").(string)

	// Executing an HTTP request
	body, err := HttpGetRequest(endpoint)
	if err != nil {
		return []GenericPowerEntry{}, err
	}

	pocketResponse, err := _unmarshalJsonToPocketbaseResponse(body)
	if err != nil {
		return []GenericPowerEntry{}, err
	}

	queryResponse := []GenericPowerEntry{}

	for i := range pocketResponse.items {
		entry := GenericPowerEntry{
			datetimeUTC: pocketResponse.items[i].HourUTC,
			priceEUR:    pocketResponse.items[i].SpotPriceEUR,
			priceDKK:    pocketResponse.items[i].SpotPriceDKK,
		}
		queryResponse = append(queryResponse, entry)
	}

	return queryResponse, nil
}

// GetQuery executes a get query towards the database it is configured with, and the given options /
// parameters provided.
//
// GetQuery() is a method on the DbConnector struct.
//
// Currently, "Pocketbase" is the only valid database connection.
//
// Returns a response consiting of a slice of GenericPowerEntry.
func (dbc DbConnector) GetQuery() ([]GenericPowerEntry, error) {
	// Check that the connector is configured as get
	if dbc.queryType != "GET" {
		return []GenericPowerEntry{}, errors.New("error: GetQuery() can only be executed on connectors configured as Get")
	}

	// Changes depending on what database we are querying
	switch dbc.database {
	case "Pocketbase":
		// After identifying Pocketbase, we now assert that the queryOptions has been configured correctly.
		switch dbc.queryOptions.(type) {
		case PocketbaseGetQueryOptions:
			// If all is good, we execute the get query and return the response
			res, err := dbc._pocketbaseGetQuery(dbc.queryOptions)
			if err != nil {
				return []GenericPowerEntry{}, errors.New(fmt.Sprintf(`Error executing get query on Pocketbase databse, got: %v`, err))
			}
			return res, nil
		default:
			return []GenericPowerEntry{}, errors.New("error: Query options does not conform to PocketbaseGetQueryOptions. See: https://app.gitbook.com/o/P816f2Z2kPDJdQmEtF9C/s/feBqL0W0wcofeJfeIX7V/backend/documentation/module-database_connector/types/pocketbasegetqueryoptions")
		}
	}
	return []GenericPowerEntry{}, errors.New("unable to determine database")
}
