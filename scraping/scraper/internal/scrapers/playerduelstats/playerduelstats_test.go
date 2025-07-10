package playerduelstats

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func getNodes() (*goquery.Selection, *goquery.Selection, *goquery.Selection, error) {
	res, err := http.Get(
		"https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf/?tab=performance",
	)
	if err != nil {
		return nil, nil, nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, nil, err
	}

	playerDuelKillsNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-normal > tbody > tr:nth-child(4) > td:nth-child(3) > div",
	)
	playerDuelFirstKillsNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-fkfd > tbody > tr:nth-child(4) > td:nth-child(3) > div",
	)
	playerDuelOpKillsNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-op> tbody > tr:nth-child(4) > td:nth-child(3) > div",
	)

	return playerDuelKillsNode, playerDuelFirstKillsNode, playerDuelOpKillsNode, nil
}

func TestPlayerDuelStatFullInfo(t *testing.T) {
	playerDuelKillsNode, playerDuelFirstKillsNode, playerDuelOpKillsNode, err := getNodes()
	if err != nil {
		t.Fatal(err)
	}

	playerDuelStatScraper := NewPlayerDuelStatScraper(
		nil,
		playerDuelKillsNode,
		playerDuelFirstKillsNode,
		playerDuelOpKillsNode,
		498628,
		5,
		17086,
		4,
	)

	if err := playerDuelStatScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	if err := playerDuelStatScraper.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
