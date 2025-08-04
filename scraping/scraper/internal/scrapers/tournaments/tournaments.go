package tournaments

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/gorm"
)

const (
	tournamentGroupSelector = `#wrapper > div.col-container > div > div.wf-card.mod-event.mod-header.mod-full > div.event-header > div.event-desc > div > div:nth-child(1) > a`
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

func tierParser(prizePool int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		return regexp.MustCompile(`vct-20[0-9][0-9]`).MatchString(strings.TrimSpace(rawVal)) ||
			prizePool >= 500000, nil
	}
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

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	tournamentSchema, ok := ctx.Value("tournamentSchema").(*models.TournamentSchema)
	if !ok {
		return fmt.Errorf("Unable to find the tournament schema")
	}

	// _, ok = ctx.Value("tx").(*gorm.DB)
	// if !ok {
	// 	return fmt.Errorf("Unable to find gorm transaction")
	// }

	if err := htmlx.ParseFromSelection(tournamentSchema, selection, htmlx.SetParsers(map[string]htmlx.Parser{
		"moneyParser": moneyParser,
	})); err != nil {
		return err
	}

	tournamentGroup, _ := selection.Find(tournamentGroupSelector).Attr("href")

	tournamentSchema.Tier1 = regexp.MustCompile(`vct-20[0-9][0-9]`).MatchString(strings.TrimSpace(tournamentGroup)) ||
		tournamentSchema.PrizePool >= 500000

	if err := helpers.PrettyPrintStruct(tournamentSchema); err != nil {
		return err
	}

	return nil
}
