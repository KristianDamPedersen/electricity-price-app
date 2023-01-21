package dataGatherer

import (
	"testing"
)

func TestCallEnergiDataService(t *testing.T) {
	validInput := []string{"2022-12-12T00:00:00+00:00", "2022-12-13T00:00:00+00:00", "DK1"}
	var err error

	// Test that it returns non-empty on valid input:
	msg, _ := CallEnergiDataService(validInput[0], validInput[1], validInput[2])
	if len(msg) == 0 {
		t.Fatalf("CallEnergiDataService() returns empty on otherwise valid input")
	}

	// Test that it throws error on invalid start date
	_, err = CallEnergiDataService("bla bla", validInput[1], validInput[2])
	if err == nil {
		t.Fatalf("CallEnergiDataService() did not return error on invalid start date")
	}

	// Test that it throws error on invalid end date
	_, err = CallEnergiDataService(validInput[0], "bla bla", validInput[2])
	if err == nil {
		t.Fatalf("CallEnergiDataService() did not return error on invalid end date")
	}

	// Test that it throws error on invalid price area
	_, err = CallEnergiDataService(validInput[0], validInput[1], "Narnia")
	if err == nil {
		t.Fatalf("CallEnergiDataService() did not return error on invalid price area")
	}

	// Test that it returns expected output (based on the valid input variable above)
	expected := PowerPriceEntry{
		HourUTC:      "2022-12-12T22:00:00",
		HourDK:       "2022-12-12T23:00:00",
		PriceArea:    "DK1",
		SpotPriceDKK: 2329.439941,
		SpotPriceEUR: 313.200012,
	}
	if msg[0] != expected {
		t.Fatalf("CallEnergiDataService() Does return the expected output, expected %v, got %v", expected, msg[0])
	}
}
