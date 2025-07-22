package tournaments

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
)

func TestTournamentScraper(t *testing.T) {
	testTournaments := []models.TournamentSchema{
		{
			Id:        2470,
			Name:      "Challengers 2025: MENA Resilience GCC-Pakistan-Iraq Split 2",
			Url:       "https://www.vlr.gg/event/2470/challengers-2025-mena-resilience-gcc-pakistan-iraq-split-2",
			PrizePool: 17500,
			Tier1:     false,
		},
		{
			Id:        2572,
			Name:      "Valorant Indonesia 2025: Summer Protocol Campus Stage",
			Url:       "https://www.vlr.gg/event/2572/valorant-indonesia-2025-summer-protocol-campus-stage",
			PrizePool: 0,
			Tier1:     false,
		},
		{
			Id:        2561,
			Name:      "EPIC.LAN #45",
			Url:       "https://www.vlr.gg/event/2561/epic-lan-45",
			PrizePool: 3017,
			Tier1:     false,
		},
		{
			Id:        2449,
			Name:      "Esports World Cup 2025",
			Url:       "https://www.vlr.gg/event/2449/esports-world-cup-2025",
			PrizePool: 1250000,
			Tier1:     true,
		},
		{
			Id:        2282,
			Name:      "Valorant Masters Toronto 2025",
			Url:       "https://www.vlr.gg/event/2282/valorant-masters-toronto-2025",
			PrizePool: 1000000,
			Tier1:     true,
		},
	}

	for _, testTournament := range testTournaments {
		res, err := http.Get(testTournament.Url)
		if err != nil {
			t.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		s := NewScraper(nil, doc.Selection, testTournament.Id, testTournament.Url)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(s.Data, testTournament); err != nil {
			t.Error(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
