package playerstats

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func TestPlayerStat(t *testing.T) {
	if err := godotenv.Load("/home/leminhohoho/repos/vlr-prediction/scraping/scraper/.env"); err != nil {
		t.Fatal(err)
	}

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
	sc.Handle(regexp.MustCompile(`playerStats`), Handler)

	res, err := http.Get(
		"https://www.vlr.gg/490310/paper-rex-vs-gen-g-champions-tour-2025-masters-toronto-r2-1-0",
	)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	selectors := []string{
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(1)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(2)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(3)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(4)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(5)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(1)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(2)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(3)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(4)",
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(5)",
	}

	for i, selector := range selectors {
		if 0 <= i && i <= 4 {
			data := Data{
				DefStat: models.PlayerOverviewStatSchema{
					Side: models.Def,
				},
				AtkStat: models.PlayerOverviewStatSchema{
					Side: models.Atk,
				},
				BothSideStat: models.PlayerOverviewStatSchema{
					Side: models.Side(""),
				},
				TeamDefRounds: 4,
				TeamAtkRounds: 12,
			}
			ctx := context.WithValue(context.Background(), "data", &data)
			ctx2 := context.WithValue(ctx, "tx", tx)

			if err := sc.Pipe("playerStats", ctx2, doc.Selection.Find(selector)); err != nil {
				t.Fatal(err)
			}
		}
		if 5 <= i && i <= 9 {
			data := Data{
				DefStat: models.PlayerOverviewStatSchema{
					Side: models.Def,
				},
				AtkStat: models.PlayerOverviewStatSchema{
					Side: models.Atk,
				},
				BothSideStat: models.PlayerOverviewStatSchema{
					Side: models.Side(""),
				},
				TeamDefRounds: 12,
				TeamAtkRounds: 4,
			}
			ctx := context.WithValue(context.Background(), "data", &data)
			ctx2 := context.WithValue(ctx, "tx", tx)

			if err := sc.Pipe("playerStats", ctx2, doc.Selection.Find(selector)); err != nil {
				t.Fatal(err)
			}
		}
	}
}
