package players

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/geographyinfo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlayerScraper struct {
	Data              models.PlayerSchema
	PlayerPageContent *goquery.Selection
	Tx                *gorm.DB
}

func NewPlayerScraper(
	tx *gorm.DB,
	playerPageContent *goquery.Selection,
	playerId int,
	playerUrl string,
) *PlayerScraper {
	return &PlayerScraper{
		Data: models.PlayerSchema{
			Id:  playerId,
			Url: playerUrl,
		},
		PlayerPageContent: playerPageContent,
		Tx:                tx,
	}
}

func (p *PlayerScraper) getCountryId(countryName string) (int, error) {
	var countryRow models.CountrySchema

	result := p.Tx.Table("countries").Where("name = ?", countryName).First(&countryRow)
	if result.Error != nil {
		return -1, result.Error
	}

	return countryRow.Id, nil
}

func (p *PlayerScraper) addCountryInfo(countryOfficialName, regionOfficialName string) (int, error) {
	regionRow := models.RegionSchema{Name: regionOfficialName}
	regionResult := p.Tx.Table("regions").Create(&regionRow)
	if regionResult.Error != nil {
		if !strings.Contains(regionResult.Error.Error(), "UNIQUE constraint failed") {
			p.Tx.Rollback()
			return -1, regionResult.Error
		}

		regionResult = p.Tx.Table("regions").Where("name = ?", regionOfficialName).First(&regionRow)
		if regionResult.Error != nil {
			return -1, regionResult.Error
		}
	}

	countryRow := models.CountrySchema{Name: countryOfficialName, RegionId: regionRow.Id}
	countryResult := p.Tx.Table("countries").Create(&countryRow)
	if countryResult.Error != nil {
		p.Tx.Rollback()
		return -1, countryResult.Error
	}

	return countryRow.Id, nil
}

func (p *PlayerScraper) countryIdParser(rawVal string) (any, error) {
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

	countryId, err := p.getCountryId(countryOfficialName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		countryId, err := p.addCountryInfo(countryOfficialName, regionOfficialName)
		if err != nil {
			return nil, err
		}

		return &countryId, nil
	}

	return &countryId, nil
}

func (p *PlayerScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(p.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (p *PlayerScraper) Scrape() error {
	parsers := map[string]htmlx.Parser{
		"countryIdParser": p.countryIdParser,
	}

	logrus.Debug("Scraping player information")
	if err := htmlx.ParseFromSelection(&p.Data, p.PlayerPageContent, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	return nil
}
