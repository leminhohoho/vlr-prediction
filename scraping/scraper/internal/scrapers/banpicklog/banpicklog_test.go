package banpicklog

import (
	"testing"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"

	_ "github.com/mattn/go-sqlite3"
)

func TestBanPickLog(t *testing.T) {
	conn, err := helpers.GetConn("/home/leminhohoho/repos/vlr-prediction/database/vlr.db")
	if err != nil {
		t.Fatal(err)
	}

	b := NewBanPickLogScraper(
		conn,
		nil,
		498628,
		624,
		2593,
		"PRX",
		"FNC",
		"PRX ban Haven; PRX ban Ascent; PRX pick Sunset; FNC pick Icebox; PRX pick Pearl; FNC pick Lotus; Split remains",
	)

	if err := b.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := b.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
