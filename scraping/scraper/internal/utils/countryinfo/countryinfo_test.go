package countryinfo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestCountryInfo(t *testing.T) {
	if err := godotenv.Load("/home/leminhohoho/repos/vlr-prediction/scraping/scraper/.env"); err != nil {
		t.Fatal(err)
	}

	countries := []string{"Singapore", "RUSSIA", "canada", "United States", "Vietnam"}

	for _, country := range countries {
		countryInfo, err := GetCountryInfo(country)
		if err != nil {
			t.Fatal(err)
		}

		jsonDat, err := json.MarshalIndent(countryInfo, "", "	")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(jsonDat))
	}
}
