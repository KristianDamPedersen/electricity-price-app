package databaseConnector

import (
	"fmt"
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
func TestPostQuery(t *testing.T) {
	// Valid configuration
	dbc := DbConnector{
		queryType:    "POST",
		database:     "Pocketbase",
		queryOptions: PocketbasePostQueryOptions{},
	}

	fmt.Println(dbc)

	sampleQueryOptions := PocketbasePostQueryOptions{
		datetimeUTC: "2022-01-01 10:00:00.123Z",
		datetimeDK:  "2022-01-01 10:00:00.123Z",
		priceArea:   "DK1",
		priceDKK:    123.123,
		priceEUR:    123.123,
	}

	// Test that it returns error if it is not configured as a POST request
	dbcInvalidType := DbConnector{
		queryType:    "INVALID TYPE",
		database:     "Pocketbase",
		queryOptions: sampleQueryOptions,
	}
	_, err := dbcInvalidType.PostQuery()
	if err == nil {
		t.Fatalf("PostQuery() does not return error on invalid type")
	}

	// Test that it returns error if it is not configured with a valid database
	dbcInvalidDatabase := DbConnector{
		queryType:    "POST",
		database:     "INVALID DATABASE",
		queryOptions: sampleQueryOptions,
	}
	_, err = dbcInvalidDatabase.PostQuery()
	if err == nil {
		t.Fatalf("PostQuery() does not return error on invalid database")
	}

	// Test that it returns error when querying Pocketbase with invalid query options
	dbcInvalidQueryOptions := DbConnector{
		queryType: "POST",
		database:  "Pocketbase",
		queryOptions: PocketbaseGetQueryOptions{
			page:    1,
			perPage: 30,
			sort:    "",
			filter:  "",
			expand:  "",
		},
	}
	_, err = dbcInvalidQueryOptions.PostQuery()
	if err == nil {
		t.Fatalf("PostQuery() does not return error on invalid query options when querying Pocketbase")
	}

	// Test that it throws an error if it recieves unexpected reply from
}

// ### UPDATE QUERY

// ### DELETE QUERY
