package countryinfo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCountryInfo(t *testing.T) {
	countries := []string{"Singapore", "RUSSIA", "canada", "USA", "Vietnam"}

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
