package matchmaps

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func intPtr(num int) *int {
	return &num
}

func TestMatchMapScraper(t *testing.T) {
	testMaps := []models.MatchMapSchema{
		{
			MatchId:       510154,
			Team1Id:       17,
			Team2Id:       1001,
			MapId:         1,
			Duration:      intPtr(3493),
			Team1DefScore: 5,
			Team1AtkScore: 2,
			Team1OTScore:  0,
			Team2DefScore: 6,
			Team2AtkScore: 7,
			Team2OTScore:  0,
			TeamDefFirst:  17,
			TeamPick:      intPtr(1001),
		},
		{
			MatchId:       510154,
			Team1Id:       17,
			Team2Id:       1001,
			MapId:         4,
			Duration:      intPtr(2935),
			Team1DefScore: 7,
			Team1AtkScore: 6,
			Team1OTScore:  0,
			Team2DefScore: 3,
			Team2AtkScore: 5,
			Team2OTScore:  0,
			TeamDefFirst:  17,
			TeamPick:      intPtr(17),
		},
		{
			MatchId:       510154,
			Team1Id:       17,
			Team2Id:       1001,
			MapId:         6,
			Duration:      intPtr(4287),
			Team1DefScore: 6,
			Team1AtkScore: 6,
			Team1OTScore:  2,
			Team2DefScore: 6,
			Team2AtkScore: 6,
			Team2OTScore:  4,
			TeamDefFirst:  17,
			TeamPick:      nil,
		},
	}
	testGamesId := []int{225041, 225042, 225043}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()

	res, err := http.Get("https://www.vlr.gg/510154/gen-g-vs-team-heretics-esports-world-cup-2025-sf")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	for i, testMap := range testMaps {
		mapNode := doc.Find(
			fmt.Sprintf(
				`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div[data-game-id="%d"]`,
				testGamesId[i],
			),
		)

		m := NewMatchMapScraper(tx, mapNode, testMap.MatchId, testMap.Team1Id, testMap.Team2Id)

		if err := m.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(testMap, m.Data); err != nil {
			t.Error(err)
		}

		if err := m.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
