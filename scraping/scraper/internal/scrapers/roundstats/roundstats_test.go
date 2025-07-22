package roundstats

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
)

func TestRoundStats(t *testing.T) {
	testRounds := []models.RoundStatSchema{
		{
			Team1Id: 2593,
			Team2Id: 8877,
			RoundOverviewSchema: models.RoundOverviewSchema{
				RoundNo:   1,
				TeamWon:   8877,
				WonMethod: models.Eliminate,
			},
			RoundEconomySchema: models.RoundEconomySchema{
				TeamDef:      2593,
				Team1BuyType: models.Pistol,
				Team2BuyType: models.Pistol,
				Team1Bank:    200,
				Team2Bank:    200,
			},
		},
		{
			Team1Id: 2593,
			Team2Id: 8877,
			RoundOverviewSchema: models.RoundOverviewSchema{
				RoundNo:   2,
				TeamWon:   8877,
				WonMethod: models.Eliminate,
			},
			RoundEconomySchema: models.RoundEconomySchema{
				TeamDef:      2593,
				Team1BuyType: models.Eco,
				Team2BuyType: models.SemiBuy,
				Team1Bank:    7300,
				Team2Bank:    6300,
			},
		},
		{
			Team1Id: 2593,
			Team2Id: 8877,
			RoundOverviewSchema: models.RoundOverviewSchema{
				RoundNo:   3,
				TeamWon:   8877,
				WonMethod: models.SpikeExplode,
			},
			RoundEconomySchema: models.RoundEconomySchema{
				TeamDef:      2593,
				Team1BuyType: models.SemiBuy,
				Team2BuyType: models.SemiBuy,
				Team1Bank:    800,
				Team2Bank:    12400,
			},
		},
		{
			Team1Id: 2593,
			Team2Id: 8877,
			RoundOverviewSchema: models.RoundOverviewSchema{
				RoundNo:   4,
				TeamWon:   2593,
				WonMethod: models.Defuse,
			},
			RoundEconomySchema: models.RoundEconomySchema{
				TeamDef:      2593,
				Team1BuyType: models.SemiBuy,
				Team2BuyType: models.FullBuy,
				Team1Bank:    8400,
				Team2Bank:    22300,
			},
		},
		{
			Team1Id: 2593,
			Team2Id: 8877,
			RoundOverviewSchema: models.RoundOverviewSchema{
				RoundNo:   5,
				TeamWon:   2593,
				WonMethod: models.Eliminate,
			},
			RoundEconomySchema: models.RoundEconomySchema{
				TeamDef:      2593,
				Team1BuyType: models.FullBuy,
				Team2BuyType: models.FullBuy,
				Team1Bank:    15200,
				Team2Bank:    12600,
			},
		},
	}

	overviewRes, err := http.Get("https://www.vlr.gg/510149/fnatic-vs-karmine-corp-esports-world-cup-2025-qf")
	if err != nil {
		t.Fatal(err)
	}

	overviewDoc, err := goquery.NewDocumentFromReader(overviewRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	economyRes, err := http.Get(
		"https://www.vlr.gg/510149/fnatic-vs-karmine-corp-esports-world-cup-2025-qf/?tab=economy",
	)
	if err != nil {
		t.Fatal(err)
	}

	economyDoc, err := goquery.NewDocumentFromReader(economyRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	for i, testRound := range testRounds {
		overviewNode := overviewDoc.Find(
			fmt.Sprintf(
				"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='225026'] > div:nth-child(2) > div > div > div:nth-child(1) > div:nth-child(%d)",
				i+2,
			),
		)

		economyNode := economyDoc.Find(
			fmt.Sprintf(
				"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='225026'] > div:nth-child(3) > table > tbody > tr:nth-child(1) > td:nth-child(%d)",
				i+2,
			),
		)

		s := NewScraper(
			nil,
			overviewNode,
			economyNode,
			testRound.MatchId,
			testRound.MapId,
			testRound.Team1Id,
			testRound.Team2Id,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(testRound, s.Data); err != nil {
			t.Error(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
