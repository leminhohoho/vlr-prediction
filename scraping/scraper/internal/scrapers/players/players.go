package players

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/geographyinfo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getCountryId(tx *gorm.DB, countryName string) (int, error) {
	var countryRow models.CountrySchema

	result := tx.Table("countries").Where("name = ?", countryName).First(&countryRow)
	if result.Error != nil {
		return -1, result.Error
	}

	return countryRow.Id, nil
}

func addCountryInfo(tx *gorm.DB, countryOfficialName, regionOfficialName string) (int, error) {
	regionRow := models.RegionSchema{Name: regionOfficialName}
	regionResult := tx.Table("regions").Create(&regionRow)
	if regionResult.Error != nil {
		if !strings.Contains(regionResult.Error.Error(), "UNIQUE constraint failed") {
			tx.Rollback()
			return -1, regionResult.Error
		}

		regionResult = tx.Table("regions").Where("name = ?", regionOfficialName).First(&regionRow)
		if regionResult.Error != nil {
			return -1, regionResult.Error
		}
	}

	countryRow := models.CountrySchema{Name: countryOfficialName, RegionId: regionRow.Id}
	countryResult := tx.Table("countries").Create(&countryRow)
	if countryResult.Error != nil {
		tx.Rollback()
		return -1, countryResult.Error
	}

	return countryRow.Id, nil
}
func countryIdParser(tx *gorm.DB) htmlx.Parser {
	return func(rawVal string) (any, error) {
		countryName := strings.TrimSpace(rawVal)

		if countryName == "" {
			return nil, nil
		}

		countryInfo, err := geographyinfo.GetInfoFromCountryName(countryName)
		if err != nil {
			return -1, err
		}

		countryOfficialName := countryInfo.Country
		regionOfficialName := countryInfo.Region

		countryId, err := getCountryId(tx, countryOfficialName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			countryId, err := addCountryInfo(tx, countryOfficialName, regionOfficialName)
			if err != nil {
				return nil, err
			}

			return &countryId, nil
		}

		return &countryId, nil
	}
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	p, ok := ctx.Value("player").(*models.PlayerSchema)
	if !ok {
		return fmt.Errorf("Unable to find player schema")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find the transaction")
	}

	logrus.Debug("Scraping player information")
	if err := htmlx.ParseFromSelection(p, selection, htmlx.SetParsers(map[string]htmlx.Parser{
		"countryIdParser": countryIdParser(tx),
	})); err != nil {
		return err
	}

	jsonDat, err := json.MarshalIndent(*p, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonDat))

	logrus.Debug("Saving player to db")
	if err := tx.Table("players").Create(p).Error; err != nil {
		return err
	}

	return nil
}
