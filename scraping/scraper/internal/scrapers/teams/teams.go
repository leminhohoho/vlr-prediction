package teams

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/repos/countryrepo"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/repos/regionrepo"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/geographyinfo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	teamLocSel = "#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.team-header > div.team-header-desc > div > div.team-header-country"
)

func getRegionInfo(tx *gorm.DB, teamLoc string) (*int, error) {
	geoInfo, err := geographyinfo.GetInfoFromRegionName(teamLoc)
	if err != nil {
		return nil, err
	}

	regionRepo := regionrepo.NewRegionRepo(tx)

	regionInfo, err := regionRepo.GetRegionByName(geoInfo.Region)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		regionInfo, err = regionRepo.InsertRegion(geoInfo.Region)
		if err != nil {
			return nil, err
		}
	}

	return &regionInfo.Id, nil
}

func getCountryInfo(tx *gorm.DB, teamLoc string) (*int, error) {
	geoInfo, err := geographyinfo.GetInfoFromCountryName(teamLoc)
	if err != nil {
		return nil, err
	}

	countryOfficialName := geoInfo.Country
	regionOfficialName := geoInfo.Region

	countryRepo := countryrepo.NewCountryRepo(tx)
	regionRepo := regionrepo.NewRegionRepo(tx)

	regionInfo, err := regionRepo.GetRegionByName(regionOfficialName)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		regionInfo, err = regionRepo.InsertRegion(regionOfficialName)
		if err != nil {
			return nil, err
		}
	}

	countryInfo, err := countryRepo.GetCountryByName(countryOfficialName)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		countryInfo, err = countryRepo.InsertCountry(regionInfo.Id, countryOfficialName)
		if err != nil {
			return nil, err
		}
	}

	return &countryInfo.Id, nil
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	var err error

	teamSchema, ok := ctx.Value("teamSchema").(*models.TeamSchema)
	if !ok {
		return fmt.Errorf("Unable to find team schema")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find the transaction")
	}

	if err := htmlx.ParseFromSelection(teamSchema, selection); err != nil {
		return err
	}

	if (*teamSchema.ShorthandName == "" || teamSchema.ShorthandName == nil) &&
		!strings.Contains(teamSchema.Name, " ") &&
		len(teamSchema.Name) <= 4 {
		teamSchema.ShorthandName = &teamSchema.Name
	}

	teamLoc := strings.TrimSpace(selection.Find(teamLocSel).Clone().Children().Remove().End().Text())

	if teamSchema.CountryId, err = getCountryInfo(tx, teamLoc); err != nil && err != geographyinfo.ErrNotFound {
		return err
	}

	if teamSchema.RegionId, err = getRegionInfo(tx, teamLoc); err != nil && err != geographyinfo.ErrNotFound {
		return err
	}

	if err := helpers.PrettyPrintStruct(teamSchema); err != nil {
		return err
	}

	logrus.Debug("Saving team info to db")
	if err := tx.Table("teams").Create(teamSchema).Error; err != nil {
		return err
	}

	return nil
}
