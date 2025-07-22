package geographyinfo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
)

var (
	ErrNotFound = errors.New("Can't find matching geography info for given location")
)

type GeographyInfo struct {
	Country string `csv:"country"`
	Region  string `csv:"continent"`
}

func getCountriesData() ([]*GeographyInfo, error) {
	var countriesInfo []*GeographyInfo
	countriesDbPath := os.Getenv("COUNTRIES_DB_PATH")

	if countriesDbPath == "" {
		return countriesInfo, fmt.Errorf("Can't find the countries database")
	}

	countryDb, err := os.Open(countriesDbPath)
	if err != nil {
		return countriesInfo, err
	}

	defer countryDb.Close()

	if err := gocsv.UnmarshalFile(countryDb, &countriesInfo); err != nil {
		return countriesInfo, err
	}

	return countriesInfo, nil
}

func GetInfoFromCountryName(countryName string) (GeographyInfo, error) {
	countriesInfo, err := getCountriesData()
	if err != nil {
		return GeographyInfo{}, err
	}

	for _, countryInfo := range countriesInfo {
		countryOfficialName := helpers.ToSnakeCase(countryInfo.Country)
		if strings.Contains(countryOfficialName, helpers.ToSnakeCase(countryName)) {
			return *countryInfo, nil
		}
	}

	return GeographyInfo{}, ErrNotFound
}

func GetInfoFromRegionName(regionName string) (GeographyInfo, error) {
	gepgraphyInfos, err := getCountriesData()
	if err != nil {
		return GeographyInfo{}, err
	}

	for _, geographyInfo := range gepgraphyInfos {
		regionOfficialName := helpers.ToSnakeCase(geographyInfo.Region)
		if strings.Contains(regionOfficialName, helpers.ToSnakeCase(regionName)) {
			return *geographyInfo, nil
		}
	}

	return GeographyInfo{}, ErrNotFound
}
