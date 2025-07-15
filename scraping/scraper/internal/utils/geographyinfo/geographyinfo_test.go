package geographyinfo

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
		geogrphyInfo, err := GetInfoFromCountryName(country)
		if err != nil {
			t.Fatal(err)
		}

		jsonDat, err := json.MarshalIndent(geogrphyInfo, "", "	")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(jsonDat))
	}

	regions := []string{"Asia", "North AMERICA", "europe"}
	for _, region := range regions {
		geogrphyInfo, err := GetInfoFromRegionName(region)
		if err != nil {
			t.Fatal(err)
		}

		jsonDat, err := json.MarshalIndent(geogrphyInfo, "", "	")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(jsonDat))
	}
}
