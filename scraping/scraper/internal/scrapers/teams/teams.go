package teams

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/repos/countryrepo"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/repos/regionrepo"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/geographyinfo"
	"gorm.io/gorm"
)

const (
	teamLocSel = "#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.team-header > div.team-header-desc > div > div.team-header-country"
)

type TeamScraper struct {
	Data            models.TeamSchema
	TeamPageContent *goquery.Selection
	Tx              *gorm.DB
}

func NewScraper(tx *gorm.DB, teamPageContent *goquery.Selection, id int, url string) *TeamScraper {
	return &TeamScraper{
		Data: models.TeamSchema{
			Id:  id,
			Url: url,
		},
		TeamPageContent: teamPageContent,
		Tx:              tx,
	}
}

func (t *TeamScraper) getRegionInfo(teamLoc string) error {
	geoInfo, err := geographyinfo.GetInfoFromRegionName(teamLoc)
	if err != nil {
		return err
	}

	regionRepo := regionrepo.NewRegionRepo(t.Tx)

	regionInfo, err := regionRepo.GetRegionByName(geoInfo.Region)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		regionInfo, err = regionRepo.InsertRegion(geoInfo.Region)
		if err != nil {
			return err
		}
	}

	t.Data.RegionId = &regionInfo.Id

	return nil
}

func (t *TeamScraper) getCountryInfo(teamLoc string) error {
	geoInfo, err := geographyinfo.GetInfoFromCountryName(teamLoc)
	if err != nil {
		return err
	}

	countryOfficialName := geoInfo.Country
	regionOfficialName := geoInfo.Region

	countryRepo := countryrepo.NewCountryRepo(t.Tx)
	regionRepo := regionrepo.NewRegionRepo(t.Tx)

	regionInfo, err := regionRepo.GetRegionByName(regionOfficialName)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		regionInfo, err = regionRepo.InsertRegion(regionOfficialName)
		if err != nil {
			return err
		}
	}
	fmt.Println(regionInfo)

	countryInfo, err := countryRepo.GetCountryByName(countryOfficialName)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		countryInfo, err = countryRepo.InsertCountry(regionInfo.Id, countryOfficialName)
		if err != nil {
			return err
		}
	}
	fmt.Println(countryInfo)

	t.Data.CountryId = &countryInfo.Id

	fmt.Println(countryOfficialName)
	fmt.Println(regionOfficialName)

	return nil
}

func (t *TeamScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(t.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (t *TeamScraper) Scrape() error {
	if err := htmlx.ParseFromSelection(&t.Data, t.TeamPageContent); err != nil {
		return err
	}

	teamLoc := strings.TrimSpace(t.TeamPageContent.Find(teamLocSel).Clone().Children().Remove().End().Text())

	if err := t.getCountryInfo(teamLoc); err != nil && err != geographyinfo.ErrNotFound {
		return err
	}

	if err := t.getRegionInfo(teamLoc); err != nil && err != geographyinfo.ErrNotFound {
		return err
	}

	return nil
}
