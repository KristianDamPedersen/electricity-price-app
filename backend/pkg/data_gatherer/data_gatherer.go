package dataGatherer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type EnergiDataServiceEntry struct {
	HourUTC      string
	HourDK       string
	PriceArea    string
	SpotPriceDKK float64
	SpotPriceEUR float64
}

func CallEnergiDataService(startDate string, endDate string, priceArea string) ([]EnergiDataServiceEntry, error) {

	startDateFormatted, err := time.Parse(time.RFC3339, startDate)
	startDateTrimmed := fmt.Sprintf("%vT%v", startDateFormatted.String()[:10], startDateFormatted.String()[11:16])
	if err != nil {
		return []EnergiDataServiceEntry{}, errors.New("invalid start date format (OBS. expects input of dd-mm-yyyyThh:mm:ss+hh:mm (RFC3339)")
	}

	endDateFormatted, err := time.Parse(time.RFC3339, endDate)
	endDateTrimmed := fmt.Sprintf("%vT%v", endDateFormatted.String()[:10], endDateFormatted.String()[11:16])
	if err != nil {
		return []EnergiDataServiceEntry{}, errors.New("invalid end date format (OBS. expects input of dd-mm-yyyyThh:mm:ss+hh:mm (RFC3339)")
	}

	switch priceArea {
	case "DK1":
		break
	case "DK2":
		break
	default:
		return []EnergiDataServiceEntry{}, errors.New("invalid price area entry (OBS. Valid entries are DK1 and DK2)")
	}

	// HTML get request
	query := fmt.Sprintf(`https://api.energidataservice.dk/dataset/Elspotprices?start=%v&end=%v&filter={"priceArea":"%v"}`, startDateTrimmed, endDateTrimmed, priceArea)
	fmt.Println(query)
	resp, err := http.Get(query)
	if err != nil {
		fmt.Println(err)
		return []EnergiDataServiceEntry{}, errors.New(fmt.Sprintf(`Unable to retrieve data from endpoint, got the following error: %v`, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []EnergiDataServiceEntry{}, errors.New("Problems reading json response body from API call")
	}

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

	fmt.Println(energyEntries[0])

	return energyEntries, nil
}
