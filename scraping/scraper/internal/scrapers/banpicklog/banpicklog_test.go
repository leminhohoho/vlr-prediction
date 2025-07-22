package banpicklog

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func getMapId(tx *gorm.DB, mapName string) (int, error) {
	var vlrMap models.MapSchema

	rs := tx.Table("maps").Where("name = ?", mapName).First(&vlrMap)
	if rs.Error != nil {
		return -1, rs.Error
	}

	return vlrMap.Id, nil
}

func intPtr(num int) *int {
	return &num
}

func validateTurn(
	t *testing.T,
	tx *gorm.DB,
	turn models.BanPickLogSchema,
	matchId int,
	teamId *int,
	mapName string,
	action models.VetoAction,
	order int,
) {
	if turn.MatchId != matchId {
		t.Errorf("Wrong match id, want %d, get %d", matchId, turn.MatchId)
	}

	if turn.TeamId == nil && teamId != nil {
		t.Errorf("team id shouldn't be nil")
	} else if turn.TeamId != nil && teamId == nil {
		t.Errorf("team id should be nil")
	}

	if turn.TeamId != nil && teamId != nil && *turn.TeamId != *teamId {
		t.Errorf("Wrong team id, want %d, get %d", *teamId, *turn.TeamId)
	}

	mapId, err := getMapId(tx, mapName)
	if err != nil {
		t.Fatal(err)
	}

	if turn.MapId != mapId {
		t.Errorf("Wrong map id, want %d, get %d", mapId, turn.MapId)
	}

	if turn.VetoOrder != order {
		t.Errorf("Wrong order, want %d, get %d", mapId, turn.MapId)
	}

	if turn.VetoAction != action {
		t.Errorf("Wrong veto action, want %s, get %s", action, turn.VetoAction)
	}
}

func TestBanPickLog(t *testing.T) {
	if err := godotenv.Load("/home/leminhohoho/repos/vlr-prediction/scraping/scraper/.env"); err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()

	b := NewBanPickLogScraper(
		tx,
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

	// Checking ban pick log information
	if len(b.Data.Turns) != 7 {
		t.Fatalf("Wrong number of turns, want 7, get %d", len(b.Data.Turns))
	}

	validateTurn(t, tx, b.Data.Turns[0], 498628, intPtr(624), "Haven", models.BanMap, 1)
	validateTurn(t, tx, b.Data.Turns[1], 498628, intPtr(624), "Ascent", models.BanMap, 2)
	validateTurn(t, tx, b.Data.Turns[2], 498628, intPtr(624), "Sunset", models.PickMap, 3)
	validateTurn(t, tx, b.Data.Turns[3], 498628, intPtr(2593), "Icebox", models.PickMap, 4)
	validateTurn(t, tx, b.Data.Turns[4], 498628, intPtr(624), "Pearl", models.PickMap, 5)
	validateTurn(t, tx, b.Data.Turns[5], 498628, intPtr(2593), "Lotus", models.PickMap, 6)
	validateTurn(t, tx, b.Data.Turns[6], 498628, nil, "Split", models.RemainMap, 7)

	if err := b.PrettyPrint(); err != nil {
		t.Fatal(err)
	}
}
