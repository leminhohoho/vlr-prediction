package matchmaps

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	roundOverviewSelector = `div:nth-child(2) > div > div > div.vlr-rounds-row > div:has(div.rnd-sq.mod-win)`
	roundEconomySelector  = `div:nth-child(3) > table > tbody > tr > td:has(div.rnd-sq)`
)

func defFirstParser(t1Id, t2Id int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		sideIdentifier := strings.TrimSpace(rawVal)

		switch sideIdentifier {
		case "mod-ct":
			return t1Id, nil
		case "mod-t":
			return t2Id, nil
		}

		return nil, fmt.Errorf("Can't determine which team def first from class: %s", sideIdentifier)
	}
}

func teamPickParser(t1Id, t2Id int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		classes := strings.TrimSpace(rawVal)
		if strings.Contains(classes, "mod-1") {
			return &t1Id, nil
		} else if strings.Contains(classes, "mod-2") {
			return &t2Id, nil
		}

		return nil, nil
	}
}

func durationParser(rawVal string) (any, error) {
	timeStr := strings.TrimSpace(rawVal)

	if timeStr == "" || timeStr == "-" {
		return nil, nil
	}

	duration, err := helpers.TimeToSeconds(timeStr)
	if err != nil {
		return nil, err
	}

	return &duration, nil
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	matchMapSchema, ok := ctx.Value("matchMapSchema").(*models.MatchMapSchema)
	if !ok {
		return fmt.Errorf("Unable to find match map schema")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find gorm transaction")
	}

	parsers := map[string]htmlx.Parser{
		"defFirstParser": defFirstParser(matchMapSchema.Team1Id, matchMapSchema.Team2Id),
		"durationParser": durationParser,
		"teamPickParser": teamPickParser(matchMapSchema.Team1Id, matchMapSchema.Team2Id),
		"mapIdParser":    customparsers.MapIdParser(tx),
	}

	mapOverviewNode := selection.Eq(0)
	mapPerformanceNode := selection.Eq(1)
	mapEconomyNode := selection.Eq(2)

	logrus.Debug("Parsing information from html onto match map schema")
	if err := htmlx.ParseFromSelection(matchMapSchema, mapOverviewNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	jsonDat, err := json.MarshalIndent(*matchMapSchema, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonDat))

	logrus.Debug("Saving match map to db")
	if err := tx.Table("match_maps").Create(matchMapSchema).Error; err != nil {
		return err
	}

	logrus.Debug("Scraping players overview stats")
	t1Hashmap, t2Hashmap, err := scrapePlayersStats(tx, sc, *matchMapSchema, mapOverviewNode)
	if err != nil {
		return err
	}

	logrus.Debug("Scraping rounds stats")
	if err := scrapeRoundsStats(tx, sc, *matchMapSchema, mapOverviewNode, mapEconomyNode); err != nil {
		logrus.Errorf("Error extracting round stats: %s, rounds stats of this map won't be uploaded", err.Error())
	}

	logrus.Debug("Scraping players duel stats")
	if err := scrapePlayerDuelStats(tx, sc, *matchMapSchema, mapPerformanceNode, t1Hashmap, t2Hashmap); err != nil {
		logrus.Errorf("Error extracting players duel stats: %s, rounds stats of this map won't be uploaded", err.Error())
	}

	fmt.Println(t1Hashmap)
	fmt.Println(t2Hashmap)

	return nil
}
