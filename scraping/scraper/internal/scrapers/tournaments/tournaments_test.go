package tournaments

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
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

	cache, err := piper.NewCacheDb("/tmp/vlr_cache.db")
	if err != nil {
		t.Fatal(err)
	}

	if err = cache.Validate(); err != nil && err != piper.ErrIncorrectSchema {
		t.Fatal(err)
	} else if err == piper.ErrIncorrectSchema {
		if err = cache.Setup(); err != nil {
			t.Fatal(err)
		}
	}

	backend := piper.NewPiperBackend(&http.Client{})

	sc := piper.NewScraper(backend, cache)
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/event\/[0-9]+\/[a-z0-9\/-]+$`), Handler)

	for _, testTournament := range testTournaments {
		tournamentSchema := models.TournamentSchema{Id: testTournament.Id, Url: testTournament.Url}

		if err := sc.Get(testTournament.Url, context.WithValue(context.Background(), "tournamentSchema", &tournamentSchema), nil); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(tournamentSchema, testTournament); err != nil {
			t.Error(err)
		}
	}
}
