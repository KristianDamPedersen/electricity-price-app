package databaseConnector

import (
	"testing"
)

// ### LIST / GET QUERY
func TestGetQuery(t *testing.T) {

	// Default query option
	queryOptions := PocketbaseGetQueryOptions{
		page:    1,
		perPage: 30,
		sort:    "",
		filter:  "",
		expand:  "",
	}

	// DB Connector (valid)
	dbc := DbConnector{
		getEndpoint:       "https://electricity-price-app.fly.dev/api/collections/electricity_prices/records/",
		preferredEndpoint: "Pocketbase",
		queryOptions:      queryOptions,
	}

	// Returns not empty
	res, _ := dbc.GetQuery()
	if len(res) == 0 {
		t.Fatalf("GetQuery() returns empty")
	}

	// Returns error on invalid endpoint
	dbcInvalidEndpoint := DbConnector{
		getEndpoint:       "https://electricity-price-app.fly.dev/api/collections/electricity_prices/blablabla",
		preferredEndpoint: "Pocketbase",
		queryOptions:      queryOptions,
	}
	_, err := dbcInvalidEndpoint.GetQuery()
	if err == nil {
		t.Fatalf("GetQuery() does not return error on invalid endpoint")
	}

	//
}

// ### CREATE QUERY

// ### UPDATE QUERY

// ### DELETE QUERY
