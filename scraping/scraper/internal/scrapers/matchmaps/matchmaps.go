package matchmaps

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MatchMapScraper struct {
	Data            models.MatchMapSchema
	MapOverviewNode *goquery.Selection
	Tx              *gorm.DB
}

func NewMatchMapScraper(
	tx *gorm.DB,
	mapOverviewNode *goquery.Selection,
	matchId, team1Id, team2Id int,
) *MatchMapScraper {
	return &MatchMapScraper{
		Data: models.MatchMapSchema{
			MatchId: matchId,
			Team1Id: team1Id,
			Team2Id: team2Id,
		},
		MapOverviewNode: mapOverviewNode,
		Tx:              tx,
	}
}

func (m *MatchMapScraper) defFirstParser(rawVal string) (any, error) {
	sideIdentifier := strings.TrimSpace(rawVal)

	switch sideIdentifier {
	case "mod-ct":
		return m.Data.Team1Id, nil
	case "mod-t":
		return m.Data.Team2Id, nil
	}

	return nil, fmt.Errorf("Can't determine which team def first from class: %s", sideIdentifier)
}

func (m *MatchMapScraper) teamPickParser(rawVal string) (any, error) {
	classes := strings.TrimSpace(rawVal)
	if strings.Contains(classes, "mod-1") {
		return &m.Data.Team1Id, nil
	} else if strings.Contains(classes, "mod-2") {
		return &m.Data.Team2Id, nil
	}

	return nil, nil
}

func durationParser(rawVal string) (any, error) {
	timeStr := strings.TrimSpace(rawVal)

	if timeStr == "" {
		return nil, nil
	}

	duration, err := helpers.TimeToSeconds(timeStr)
	if err != nil {
		return nil, err
	}

	return &duration, nil
}

func (m *MatchMapScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(m.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (m *MatchMapScraper) Scrape() error {
	logrus.Debug("Parsing information from map html content into match map schema")
	parsers := map[string]htmlx.Parser{
		"defFirstParser": m.defFirstParser,
		"durationParser": durationParser,
		"teamPickParser": m.teamPickParser,
		"mapIdParser":    customparsers.MapIdParser(m.Tx),
	}

	if err := htmlx.ParseFromSelection(
		&m.Data, m.MapOverviewNode,
		htmlx.SetParsers(parsers),
	); err != nil {
		return err
	}

	return nil
}
