package playerstats

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	_ "github.com/mattn/go-sqlite3"
)

func TestPlayerStat(t *testing.T) {

	dbPath := "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
	conn, err := helpers.GetConn(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

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

	forsakenOverviewStatNode := doc.Find(
		`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(5)`,
	)

	forsakenOverviewStatScraper := NewPlayerOverviewStatScraper(
		conn,
		nil,
		forsakenOverviewStatNode,
		490310,
		624,
		17,
		12,
		4,
	)

	if err := forsakenOverviewStatScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := forsakenOverviewStatScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}

	d4v4iOverviewStatNode := doc.Find(
		`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(3) > div:nth-child(4) > div:nth-child(1) > table > tbody > tr:nth-child(3)`,
	)

	d4v4iOverviewStatScraper := NewPlayerOverviewStatScraper(
		conn,
		nil,
		d4v4iOverviewStatNode,
		490310,
		624,
		17,
		12,
		4,
	)

	if err := d4v4iOverviewStatScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := d4v4iOverviewStatScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
