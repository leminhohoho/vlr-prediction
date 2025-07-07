package roundstats

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestRoundStats(t *testing.T) {
	overviewRes, err := http.Get(
		"https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf",
	)
	if err != nil {
		t.Fatal(err)
	}

	economyRes, err := http.Get(
		"https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf/?tab=economy",
	)
	if err != nil {
		t.Fatal(err)
	}

	overviewDoc, err := goquery.NewDocumentFromReader(overviewRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	economyDoc, err := goquery.NewDocumentFromReader(economyRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	roundOverviewNode := overviewDoc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game > div:nth-child(2) > div > div > div > div:nth-child(5)",
	).First()

	roundEconomyNode := economyDoc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game > div:nth-child(3) > table > tbody > tr:nth-child(1) > td:nth-child(5)",
	).First()

	roundStatScraper := NewRoundStatScraper(
		nil,
		nil,
		roundOverviewNode,
		roundEconomyNode,
		498628,
		9,
		624,
		2593,
	)

	if err := roundStatScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := roundStatScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
