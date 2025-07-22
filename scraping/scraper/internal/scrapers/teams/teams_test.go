package teams

import (
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/table"
	"github.com/joho/godotenv"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()

	for _, teamUrl := range teamUrls {
		res, err := http.Get(teamUrl)
		if err != nil {
			t.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		ts := NewScraper(tx, doc.Selection, 0, teamUrl)

		if err := ts.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := ts.PrettyPrint(); err != nil {
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
