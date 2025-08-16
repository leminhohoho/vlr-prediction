package matches

import (
	"context"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func TestMatchScraper(t *testing.T) {
	testMatches := []models.MatchSchema{
		{
			Id:           530355,
			Url:          "https://www.vlr.gg/530355/rex-regum-qeon-vs-t1-vct-2025-pacific-stage-2-ur1",
			Date:         time.Time{},
			TournamentId: 2500,
			Stage:        models.Playoff,
			Team1Id:      878,
			Team2Id:      14,
			Team1Score:   1,
			Team2Score:   2,
			Team1Rating:  1784,
			Team2Rating:  1823,
		},
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
			Team1Rating:  1654,
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

	logrus.SetLevel(logrus.TraceLevel)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()

	cache, err := piper.NewCacheDb("/tmp/vlr_cache.db")
	if err != nil {
		t.Fatal(err)
	}

	if err = cache.Validate(); err != nil && err != piper.ErrIncorrectSchema {
		t.Fatal(err)
	} else if err == piper.ErrIncorrectSchema {
		if err = cache.Setup(); err != nil {
			t.Fatal(err)
		}
	}

	backend := piper.NewPiperBackend(&http.Client{})

	sc := piper.NewScraper(backend, cache)
	sc.Handle(regexp.MustCompile(`match`), Handler)

	for _, testMatch := range testMatches {
		res, err := http.Get(testMatch.Url)
		if err != nil {
			t.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		m := models.MatchSchema{
			Id:   testMatch.Id,
			Url:  testMatch.Url,
			Date: time.Time{},
		}

		ctx := context.WithValue(context.Background(), "matchSchema", &m)
		ctx2 := context.WithValue(ctx, "tx", tx)

		if err := sc.Pipe("match", ctx2, doc.Selection); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(testMatch, m); err != nil {
			t.Error(err)
		}
	}

	tx.Rollback()
}
