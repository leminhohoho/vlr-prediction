package matches

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Stage string

const (
	GroupStage Stage = "group_stage"
	Playoff    Stage = "playoff"
	GrandFinal Stage = "grand_final"
)

type MatchSchema struct {
	Id           int
	Url          string
	Date         time.Time
	TournamentId int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-super > div:nth-child(1) > a"                                                                             source:"attr=href" parser:"idParser"`
	Stage        Stage `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-super > div:nth-child(1) > a > div > div.match-header-event-series"                    parser:"stageParser"`
	Team1Id      int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1"                                                             source:"attr=href" parser:"idParser"`
	Team2Id      int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2"                                                             source:"attr=href" parser:"idParser"`
	Team1Score   int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-winner"`
	Team2Score   int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-loser"`
	Team1Rating  int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1 > div > div.match-header-link-name-elo"                                         parser:"ratingParser"`
	Team2Rating  int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2 > div > div.match-header-link-name-elo"                                         parser:"ratingParser"`
}

type MatchScraper struct {
	Data             MatchSchema
	MatchPageContent *goquery.Selection
	Conn             *sql.DB
	Tx               *gorm.Tx
}

func NewMatchScraper(
	conn *sql.DB,
	tx *gorm.Tx,
	htmlContent *goquery.Selection,
	id int,
	url string,
	date time.Time,
) *MatchScraper {
	return &MatchScraper{
		Data: MatchSchema{
			Id:   id,
			Url:  url,
			Date: date,
		},
		MatchPageContent: htmlContent,
		Conn:             conn,
		Tx:               tx,
	}
}

func stageParser(rawVal string) (any, error) {
	matchHeader := helpers.ToSnakeCase(strings.TrimSpace(rawVal))
	if strings.Contains(matchHeader, "grand_final") {
		return GrandFinal, nil
	}

	if strings.Contains(matchHeader, "playoff") {
		return Playoff, nil
	}

	return GroupStage, nil
}

func ratingParser(rawVal string) (any, error) {
	ratingStr := strings.TrimSpace(rawVal)
	if ratingStr == "" {
		return 0, nil
	}

	if !regexp.MustCompile(`^\[[0-9]+\]$`).MatchString(ratingStr) {
		return -1, fmt.Errorf("rating string %s is not valid for converting to rating", ratingStr)
	}

	ratingStr = ratingStr[1 : len(ratingStr)-1]
	rating, _ := strconv.Atoi(ratingStr)
	return rating, nil
}

func (m *MatchScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(m.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (m *MatchScraper) Scrape() error {
	logrus.Debug("Parsing information from match html content into match schema")
	parsers := map[string]htmlx.Parser{
		"idParser":     customparsers.IdParser,
		"stageParser":  stageParser,
		"ratingParser": ratingParser,
	}

	if err := htmlx.ParseFromSelection(&m.Data, m.MatchPageContent, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	return nil
}
