package players

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/table"
	"github.com/joho/godotenv"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func TestPlayerFullInfo(t *testing.T) {
	if err := godotenv.Load("/home/leminhohoho/repos/vlr-prediction/scraping/scraper/.env"); err != nil {
		t.Fatal(err)
	}

	playerUrls := []string{
		"https://www.vlr.gg/player/13744/patmen",
		"https://www.vlr.gg/player/7378/jinggg",
		"https://www.vlr.gg/player/9801/f0rsaken",
		"https://www.vlr.gg/player/9803/d4v41",
		"https://www.vlr.gg/player/17086/something",
		"https://www.vlr.gg/player/9800/mindfreak",
		"https://www.vlr.gg/player/438/boaster",
		"https://www.vlr.gg/player/4/crashies",
		"https://www.vlr.gg/player/9554/kaajak",
		"https://www.vlr.gg/player/458/chronicle",
		"https://www.vlr.gg/player/9810/alfajer",
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()

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
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/player\/[0-9]+\/[a-z0-9]+$`), Handler)

	for _, playerUrl := range playerUrls {
		log.Printf("Scraping %s\n", playerUrl)

		playerIdStr := strings.Split(playerUrl, "/")[4]
		playerId, _ := strconv.Atoi(playerIdStr)
		p := models.PlayerSchema{
			Id:  playerId,
			Url: playerUrl,
		}

		ctx := context.WithValue(context.Background(), "player", &p)
		ctx2 := context.WithValue(ctx, "tx", tx)

		if err := sc.Get(playerUrl, ctx2, nil); err != nil {
			t.Fatal(err)
		}
	}

	var countries []models.CountrySchema
	rs := tx.Table("countries").Find(&countries)
	if rs.Error != nil {
		t.Fatal(rs.Error)
	}

	var regions []models.RegionSchema
	rs = tx.Table("regions").Find(&regions)
	if rs.Error != nil {
		t.Fatal(rs.Error)
	}

	countryTable := table.NewWriter()
	countryTable.SetOutputMirror(os.Stdout)
	countryTable.AppendHeader(table.Row{"id", "name", "region_id"})
	for _, country := range countries {
		countryTable.AppendRow(table.Row{country.Id, country.Name, country.RegionId})
	}

	regionTable := table.NewWriter()
	regionTable.SetOutputMirror(os.Stdout)
	regionTable.AppendHeader(table.Row{"id", "name"})
	for _, regions := range regions {
		regionTable.AppendRow(table.Row{regions.Id, regions.Name})
	}

	countryTable.Render()
	regionTable.Render()

	tx.Rollback()
}
