package playerhighlights

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
)

func TestPlayerHightlights(t *testing.T) {
	res, err := http.Get(
		"https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf/?tab=performance",
	)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	playerHighlightNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='221168'] > div:nth-child(2) > table > tbody > tr:nth-child(11) > td:nth-child(3) > div > div > div > div:nth-child(2)",
	)

	prxHashMap := map[string]int{
		"f0rsakeN":  9801,
		"d4v41":     9803,
		"Jinggg":    7378,
		"PatMen":    13744,
		"something": 17086,
	}

	playerHighlightScraper := NewPlayerHighlightScraper(
		nil,
		nil,
		playerHighlightNode,
		498628,
		9,
		2593,
		9554,
		models.P1v2,
		prxHashMap,
	)

	if err := playerHighlightScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := playerHighlightScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
