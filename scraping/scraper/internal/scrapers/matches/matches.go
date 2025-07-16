package matches

import (
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
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MatchScraper struct {
	Data             models.MatchSchema
	MatchPageContent *goquery.Selection
	Tx               *gorm.DB
}

func NewMatchScraper(
	tx *gorm.DB,
	htmlContent *goquery.Selection,
	id int,
	url string,
	date time.Time,
) *MatchScraper {
	return &MatchScraper{
		Data: models.MatchSchema{
			Id:   id,
			Url:  url,
			Date: date,
		},
		MatchPageContent: htmlContent,
		Tx:               tx,
	}
}

func stageParser(rawVal string) (any, error) {
	matchHeader := helpers.ToSnakeCase(strings.TrimSpace(rawVal))
	if strings.Contains(matchHeader, "grand_final") {
		return models.GrandFinal, nil
	}

	if strings.Contains(matchHeader, "playoff") {
		return models.Playoff, nil
	}

	return models.GroupStage, nil
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
