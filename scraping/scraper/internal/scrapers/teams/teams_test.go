package teams

import (
	"context"
	"net/http"
	"os"
	"regexp"
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

func TestTeamScraperWithCountryName(t *testing.T) {
	teamUrls := []string{
		"https://www.vlr.gg/team/2593/fnatic",
		"https://www.vlr.gg/team/624/paper-rex",
		"https://www.vlr.gg/team/397/bbl-esports",
		"https://www.vlr.gg/team/1034/nrg",
		"https://www.vlr.gg/team/8877/karmine-corp",
		"https://www.vlr.gg/team/2/sentinels",
		"https://www.vlr.gg/team/1001/team-heretics",
		"https://www.vlr.gg/team/17/gen-g",
	}

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable all logs
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
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/team\/[0-9]+\/[a-z0-9-]+$`), Handler)

	for _, teamUrl := range teamUrls {
		teamSchema := models.TeamSchema{Url: teamUrl}

		ctx := context.WithValue(context.Background(), "teamSchema", &teamSchema)
		ctx2 := context.WithValue(ctx, "tx", tx)

		if err := sc.Get(teamUrl, ctx2, nil); err != nil {
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
