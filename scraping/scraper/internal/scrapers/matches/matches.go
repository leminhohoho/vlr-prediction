package matches

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/urlinfo"
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
	TournamentId int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-super > div:nth-child(1) > a"                                                                             source:"attr=href" parser:"urlParser"`
	Stage        Stage `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-super > div:nth-child(1) > a > div > div.match-header-event-series"                    parser:"stageParser"`
	Team1Id      int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1"                                                             source:"attr=href" parser:"urlParser"`
	Team2Id      int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2"                                                             source:"attr=href" parser:"urlParser"`
	Team1Score   int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-winner"`
	Team2Score   int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-loser"`
	Team1Rating  int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1 > div > div.match-header-link-name-elo"                                         parser:"ratingParser"`
	Team2Rating  int   `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2 > div > div.match-header-link-name-elo"                                         parser:"ratingParser"`
}

type MatchScraper struct {
	Data        MatchSchema
	HtmlContent *goquery.Selection
	Conn        *sql.DB
	Tx          *gorm.Tx
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
		HtmlContent: htmlContent,
		Conn:        conn,
		Tx:          tx,
	}
}

func urlParser(rawVal string) (any, error) {
	url := strings.TrimSpace(rawVal)
	vlrUrlInfo, err := urlinfo.ExtractUrlInfo(url)
	if err != nil {
		return nil, err
	}

	return vlrUrlInfo.Id, nil
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

func (m *MatchScraper) PrettyPrint() {
	fmt.Printf("Match id: %d\n", m.Data.Id)
	fmt.Printf("Match url: %s\n", m.Data.Url)
	fmt.Printf("Match date: %s\n", m.Data.Date.Format("Mon, January 2, 2006"))
	fmt.Printf("Match tournament id: %d\n", m.Data.TournamentId)
	fmt.Printf("Match stage: %s\n", m.Data.Stage)
	fmt.Printf("Match team 1 id: %d\n", m.Data.Team1Id)
	fmt.Printf("Match team 2 id: %d\n", m.Data.Team2Id)
	fmt.Printf("Match team 1 score: %d\n", m.Data.Team1Score)
	fmt.Printf("Match team 2 score: %d\n", m.Data.Team2Score)
	fmt.Printf("Match team 1 rating: %d\n", m.Data.Team1Rating)
	fmt.Printf("Match team 2 rating: %d\n", m.Data.Team2Rating)
}

func (m *MatchScraper) Scrape() error {
	logrus.Debug("Parsing information from match html content into match schema")
	parsers := map[string]htmlx.Parser{
		"urlParser":    urlParser,
		"stageParser":  stageParser,
		"ratingParser": ratingParser,
	}

	if err := htmlx.ParseFromSelection(&m.Data, m.HtmlContent, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	return nil
}
