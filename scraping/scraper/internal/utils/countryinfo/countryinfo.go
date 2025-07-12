package countryinfo

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
)

type CountryInfo struct {
	Name       string `csv:"country"`
	RegionName string `csv:"continent"`
}

type CountriesInfo struct {
	Data []CountryInfo
}

func GetCountryInfo(countryName string) (CountryInfo, error) {
	var countriesInfo []*CountryInfo
	countriesDbPath := os.Getenv("COUNTRIES_DB_PATH")

	if countriesDbPath == "" {
		return CountryInfo{}, fmt.Errorf("Can't find the countries database")
	}

	countryDb, err := os.Open(countriesDbPath)
	if err != nil {
		return CountryInfo{}, err
	}
	defer countryDb.Close()

	if err := gocsv.UnmarshalFile(countryDb, &countriesInfo); err != nil {
		return CountryInfo{}, err
	}

	for _, countryInfo := range countriesInfo {
		countryOfficialName := helpers.ToSnakeCase(countryInfo.Name)
		if strings.Contains(countryOfficialName, helpers.ToSnakeCase(countryName)) {
			return *countryInfo, nil
		}
	}

	return CountryInfo{}, fmt.Errorf("Country name doesn't match any thing in the DB")
}
