package matchmaps

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestMapWithFullInfo(t *testing.T) {
	res, err := http.Get("https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	doc.Find(
		`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div[data-game-id!="all"].vm-stats-game:has(div + div)`,
	).Each(func(_ int, mapNode *goquery.Selection) {
		matchMapScraper := NewMatchMapScraper(
			nil,
			nil,
			mapNode,
			-1,
			498628,
			624,
			2593,
		)

		if err := matchMapScraper.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := matchMapScraper.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestMapWithMissingInfo(t *testing.T) {
	res, err := http.Get("https://www.vlr.gg/490314/paper-rex-vs-team-liquid-valorant-masters-toronto-2025-r3-1-1")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	doc.Find(
		`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div[data-game-id!="all"].vm-stats-game:has(div + div)`,
	).Each(func(_ int, mapNode *goquery.Selection) {
		matchMapScraper := NewMatchMapScraper(
			nil,
			nil,
			mapNode,
			-1,
			490314,
			624,
			474,
		)

		if err := matchMapScraper.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := matchMapScraper.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	})
}
