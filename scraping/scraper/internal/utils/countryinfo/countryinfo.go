package countryinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NameInfo struct {
	Name string `json:"common"`
}

type CountryInfo struct {
	Name       NameInfo `json:"name"`
	RegionName string   `json:"region"`
}

type CountriesInfo struct {
	Data []CountryInfo
}

func GetCountryInfo(countryName string) (CountryInfo, error) {
	var countryInfo CountryInfo

	res, err := http.Get("https://restcountries.com/v3.1/name/" + countryName)
	if err != nil {
		return countryInfo, err
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return countryInfo, err
	}

	var countriesInfo CountriesInfo

	if err := json.Unmarshal(dat, &countriesInfo.Data); err != nil {
		return countryInfo, err
	}

	if len(countriesInfo.Data) == 0 {
		return countryInfo, fmt.Errorf("Country info is empty")
	}

	return countriesInfo.Data[0], err
}
