package matches

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func compareMatches(t *testing.T, test models.MatchSchema, result models.MatchSchema) {
	sTest := reflect.ValueOf(test)
	sResult := reflect.ValueOf(result)

	for i := range sTest.NumField() {
		vTest := sTest.Field(i)
		vResult := sResult.Field(i)
		fieldName := sTest.Type().Field(i).Name

		if !vTest.Equal(vResult) {
			t.Errorf(
				"Error validating field '%s', want '%v', get '%v'",
				fieldName,
				vTest.Interface(),
				vResult.Interface(),
			)
		}
	}
}

func TestMatchScraper(t *testing.T) {
	testMatches := []models.MatchSchema{
		{
			Id:           506931,
			Url:          "https://www.vlr.gg/506931/bilibili-gaming-vs-tyloo-vct-2025-china-stage-2-w1",
			Date:         time.Time{},
			TournamentId: 2499,
			Stage:        models.GroupStage,
			Team1Id:      12010,
			Team2Id:      731,
			Team1Score:   2,
			Team2Score:   0,
			Team1Rating:  1719,
			Team2Rating:  1423,
		},
		{
			Id:           489315,
			Url:          "https://www.vlr.gg/489315/twisted-minds-vs-villianarc-challengers-2025-mena-resilience-gcc-pakistan-iraq-split-2-w5",
			Date:         time.Time{},
			TournamentId: 2470,
			Stage:        models.GroupStage,
			Team1Id:      6035,
			Team2Id:      18504,
			Team1Score:   2,
			Team2Score:   0,
			Team1Rating:  1808,
			Team2Rating:  0,
		},
		{
			Id:           510155,
			Url:          "https://www.vlr.gg/510155/fnatic-vs-team-heretics-esports-world-cup-2025-gf",
			Date:         time.Time{},
			TournamentId: 2449,
			Stage:        models.GrandFinal,
			Team1Id:      2593,
			Team2Id:      1001,
			Team1Score:   2,
			Team2Score:   3,
			Team1Rating:  2153,
			Team2Rating:  1974,
		},
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()

	for _, testMatch := range testMatches {
		res, err := http.Get(testMatch.Url)
		if err != nil {
			t.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		m := NewMatchScraper(tx, doc.Selection, testMatch.Id, testMatch.Url, time.Time{})

		if err := m.Scrape(); err != nil {
			t.Fatal(err)
		}

		compareMatches(t, testMatch, m.Data)

		if err := m.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
