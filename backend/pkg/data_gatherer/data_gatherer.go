package dataGatherer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Serves as intermediate conversion, we output a [] of these after calling the API through CallEnergiDataService.
//
// For documentation see: https://app.gitbook.com/o/P816f2Z2kPDJdQmEtF9C/s/feBqL0W0wcofeJfeIX7V/backend/documentation/module-data_gatherer/types/energidataserviceentry
type EnergiDataServiceEntry struct {
	HourUTC      string
	HourDK       string
	PriceArea    string
	SpotPriceDKK float64
	SpotPriceEUR float64
}

// Converts string from RFC3339 to a string compatible with our API endpoint
func _ParseRFC3339StringToAPIString(s string) (string, error) {
	str, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return "", errors.New("invalid date format. (OBS: Expects RFC3339 (dd-mm-yyyyThh:mm:ss+hh:mm)")
	}
	strTrimmed := fmt.Sprintf("%vT%v", str.String()[0:10], str.String()[11:16])
	return strTrimmed, nil
}

// Checks that the pricing area is valid and returns it if so
func _ParsePriceAreaToAPI(priceArea string) (string, error) {
	switch priceArea {
	case "DK1":
		return "DK1", nil
	case "DK2":
		return "DK2", nil
	default:
		return "", errors.New("invalid price area entry (OBS. Valid entries are DK1 and DK2)")
	}
}

// Unmarshals the response from our API call into a slice of EnergiDataServiceEntry's
func _UnmarshalToEnergiDataServiceEntries(body []byte) []EnergiDataServiceEntry {
	var jsonResponse map[string]any
	json.Unmarshal(body, &jsonResponse)

	records := jsonResponse["records"].([]interface{})

	var energyEntries []EnergiDataServiceEntry

	for i := range records {
		entry := records[i].(map[string]interface{})
		energyEntry := EnergiDataServiceEntry{
			HourUTC:      entry["HourUTC"].(string),
			HourDK:       entry["HourDK"].(string),
			PriceArea:    entry["PriceArea"].(string),
			SpotPriceDKK: entry["SpotPriceDKK"].(float64),
			SpotPriceEUR: entry["SpotPriceEUR"].(float64),
		}
		energyEntries = append(energyEntries, energyEntry)
	}

	return energyEntries
}

// Executes a query against our API endpoint
func _queryEnergiDataService(query string) ([]byte, error) {
	resp, err := http.Get(query)
	if err != nil {
		fmt.Println(err)
		return []byte{}, errors.New(fmt.Sprintf(`Unable to retrieve data from endpoint, got the following error: %v`, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Problems reading json response body from API call")
	}

	return body, nil
}

// CallEnergiDataService queries the Energi Data Service API endpoint for power prices.
//
// You can find the endpoint here: https://api.energidataservice.dk/dataset/Elspotprices
//
// For further documentation: https://app.gitbook.com/o/P816f2Z2kPDJdQmEtF9C/s/feBqL0W0wcofeJfeIX7V/backend/documentation/module-data_gatherer/callenergidataservice
//
// startDate represents the starting date that we are interested in (RFC3339 formatted)
//
// endDate represents the end date that we are interested in (RFC3339 formatted)
//
// priceArea specifies that pricing area we are interested in (currently supports "DK1" and "DK2")
//
// Returns a slice of EnergiDataServiceEntry's (these entries are spaced 1 hour apart)
func CallEnergiDataService(startDate string, endDate string, priceArea string) ([]EnergiDataServiceEntry, error) {

	startDateFormatted, err := _ParseRFC3339StringToAPIString(startDate)
	if err != nil {
		return []EnergiDataServiceEntry{}, err
	}

	endDateFormatted, err := _ParseRFC3339StringToAPIString(endDate)
	if err != nil {
		return []EnergiDataServiceEntry{}, err
	}

	priceAreaFormatted, err := _ParsePriceAreaToAPI(priceArea)
	if err != nil {
		return []EnergiDataServiceEntry{}, err
	}

	// HTML get request
	query := fmt.Sprintf(`https://api.energidataservice.dk/dataset/Elspotprices?start=%v&end=%v&filter={"priceArea":"%v"}`, startDateFormatted, endDateFormatted, priceAreaFormatted)
	resp, _ := _queryEnergiDataService(query)

	res := _UnmarshalToEnergiDataServiceEntries(resp)

	fmt.Println(res)
	return res, nil
}
