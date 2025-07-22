package tournaments

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/gorm"
)

type TournamentScraper struct {
	Data                  models.TournamentSchema
	TournamentPageContent *goquery.Selection
	Tx                    *gorm.DB
}

func NewScraper(tx *gorm.DB, tournamentPageContent *goquery.Selection, id int, url string) *TournamentScraper {
	return &TournamentScraper{
		Data: models.TournamentSchema{
			Id:  id,
			Url: url,
		},
		TournamentPageContent: tournamentPageContent,
		Tx:                    tx,
	}
}

func moneyParser(rawVal string) (any, error) {
	moneyStr := strings.TrimSpace(rawVal)
	if moneyStr == "" {
		return 0, nil
	}

	prizeStr := regexp.MustCompile(`[0-9,]+ USD`).FindString(moneyStr)
	if moneyStr == "" {
		return nil, fmt.Errorf("Error extracting amount of money from: %s", moneyStr)
	}

	prizeStr = strings.TrimSpace(strings.Replace(strings.Replace(prizeStr, "USD", "", -1), ",", "", -1))
	prize, _ := strconv.Atoi(prizeStr)

	return prize, nil
}

func (t *TournamentScraper) tierParser(rawVal string) (any, error) {
	return regexp.MustCompile(`vct-20[0-9][0-9]`).MatchString(strings.TrimSpace(rawVal)) ||
		t.Data.PrizePool >= 500000, nil
}

func (t *TournamentScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(t.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (t *TournamentScraper) Scrape() error {
	parsers := map[string]htmlx.Parser{
		"moneyParser": moneyParser,
		"tierParser":  t.tierParser,
	}

	if err := htmlx.ParseFromSelection(&t.Data, t.TournamentPageContent, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	return nil
}
