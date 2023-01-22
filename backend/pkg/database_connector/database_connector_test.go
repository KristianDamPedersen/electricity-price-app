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
		queryType:    "GET",
		database:     "Pocketbase",
		queryOptions: queryOptions,
	}

	// Returns not empty
	res, _ := dbc.GetQuery()
	if len(res) == 0 {
		t.Fatalf("GetQuery() returns empty")
	}

	// test that GetQuery returns error on invalid database
	dbcInvalidDatabase := DbConnector{
		database: "Bla bla bla",
		queryOptions: PocketbaseGetQueryOptions{
			page:    1,
			perPage: 30,
			sort:    "",
			filter:  "",
			expand:  "",
		},
	}
	_, err := dbcInvalidDatabase.GetQuery()
	if err == nil {
		t.Fatalf("GetQuery() does not return error on invalid database")
	}

	// test that GetQuery returns error on invalid PocketbaseGetQueryFormatting
	dbcInvalidPocketbaseQueryOptions := DbConnector{
		queryType:    "Get",
		database:     "Pocketbase",
		queryOptions: nil,
	}
	_, err = dbcInvalidPocketbaseQueryOptions.GetQuery()
	if err == nil {
		t.Fatalf("GetQuery() does not return error on invalid pocketbase query formatting")
	}

	// Test that GetQuery (on Pocketbase) returns error on unexpected return value ( !TODO this will break once the other queries are in place )
	expected := GenericPowerEntry{
		datetimeUTC: "2023-01-04 12:00:00.000Z",
		priceEUR:    35.01,
		priceDKK:    20.21,
	}

	res, _ = dbc.GetQuery()
	if res[0] != expected {
		t.Fatalf("Recieved wrong return value from pocketbase, expected %v but got %v", expected, res)
	}

	// Test that GetQuery (on pocketbase) returns error if its not configured as Get query.
	dbcInvalidType := DbConnector{
		queryType: "INVALID TYPE",
		database:  "Pocketbase",
		queryOptions: PocketbaseGetQueryOptions{
			page:    1,
			perPage: 30,
			sort:    "",
			filter:  "",
			expand:  "",
		},
	}
	_, err = dbcInvalidType.GetQuery()
	if err == nil {
		t.Fatalf("GetQuery() does not return error on invalid type")
	}
}

// ### CREATE QUERY

// ### UPDATE QUERY

// ### DELETE QUERY
