package players

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	_ "github.com/mattn/go-sqlite3"
)

func TestPlayerFullInfo(t *testing.T) {
	playerUrl := "https://www.vlr.gg/player/9801/f0rsaken"

	res, err := http.Get(playerUrl)
	if err != nil {
		t.Fatal(err)
	}

	dbPath := "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
	conn, err := helpers.GetConn(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	playerScraper := NewPlayerScraper(conn, nil, doc.Selection, 9801, playerUrl)
	if err := playerScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := playerScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
