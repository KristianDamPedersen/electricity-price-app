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
	getEndpoint  string
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

type PocketbaseResponse struct {
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

	// Unmarshalling into PocketbaseResponse
	fmt.Println(int64(jsonResponse["page"].(float64))) // page number
	fmt.Println(jsonResponse["perPage"])               // perPage
	fmt.Println(jsonResponse["totalItems"])            // totalItems
	fmt.Println(jsonResponse["totalPages"])            // totalPages

	fmt.Println(jsonResponse["items"].([]interface{})[0].(map[string]any)["priceDKK"])
	fmt.Println(jsonResponse["items"].([]interface{})[0].(map[string]any)["priceEUR"])
	fmt.Println(jsonResponse["items"].([]interface{})[0].(map[string]any)["datetime_UTC"])
	fmt.Println(jsonResponse["items"].([]interface{})[0].(map[string]any)["datetime_DK"])
	fmt.Println(jsonResponse["items"].([]interface{})[0].(map[string]any)["priceArea"])

	pocketResponse := PocketbaseResponse{
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

func (dbc DbConnector) GetQuery() ([]GenericPowerEntry, error) {
	switch dbc.database {
	case "Pocketbase":
		switch dbc.queryOptions.(type) {
		case PocketbaseGetQueryOptions:
			res, err := dbc._pocketbaseGetQuery(dbc.queryOptions)
			if err != nil {
				return []GenericPowerEntry{}, errors.New(fmt.Sprintf(`Error executing get query on Pocketbase databse, got: %v`, err))
			}
			return res, nil
		default:
			return []GenericPowerEntry{}, errors.New("Error: Query options does not conform to PocketbaseGetQueryOptions. See: ") // !TODO documentation here
		}
	}
	return []GenericPowerEntry{}, errors.New("Unable to determine database, please visit: xxx to see valid database options") // !TODO documentation here
}
