package playerstats

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
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
			p := NewPlayerOverviewStatScraper(tx, doc.Selection.Find(selector), -1, -1, -1, 4, 12)
			if err := p.Scrape(); err != nil {
				t.Fatal(err)
			}
			if err := p.PrettyPrint(); err != nil {
				t.Fatal(err)
			}
		}
		if 5 <= i && i <= 9 {
			p := NewPlayerOverviewStatScraper(tx, doc.Selection.Find(selector), -1, -1, -1, 12, 4)
			if err := p.Scrape(); err != nil {
				t.Fatal(err)
			}
			if err := p.PrettyPrint(); err != nil {
				t.Fatal(err)
			}
		}
	}
}
