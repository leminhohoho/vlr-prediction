package playerhighlights

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbPath = "/home/leminhohoho/repos/vlr-prediction/database/vlr.db"
)

func compareHighlights(highlightA, highlightB []models.PlayerHighlightSchema) error {
	if len(highlightA) != len(highlightB) {
		return fmt.Errorf(
			"highlight a of length %d is not the same as highlight b of length %d ",
			len(highlightA),
			len(highlightB),
		)
	}

	for _, highlightALog := range highlightA {
		if !slices.ContainsFunc(highlightB, func(highlightB models.PlayerHighlightSchema) bool {
			return highlightB.RoundNo == highlightALog.RoundNo &&
				highlightB.HighlightType == highlightALog.HighlightType &&
				highlightB.PlayerAgainstId == highlightALog.PlayerAgainstId
		}) {
			return fmt.Errorf("Unable to find match highlight log for %v", highlightALog)
		}
	}

	return nil
}

func TestPlayerHightlights(t *testing.T) {
	p2ks := [][]models.PlayerHighlightSchema{
		{
			{RoundNo: 6, HighlightType: models.P2k, PlayerAgainstId: 727},
			{RoundNo: 6, HighlightType: models.P2k, PlayerAgainstId: 16459},
		},
		{
			{RoundNo: 8, HighlightType: models.P2k, PlayerAgainstId: 727},
			{RoundNo: 8, HighlightType: models.P2k, PlayerAgainstId: 10307},
		},
		{
			{RoundNo: 19, HighlightType: models.P2k, PlayerAgainstId: 10307},
			{RoundNo: 19, HighlightType: models.P2k, PlayerAgainstId: 5654},
		},
	}
	p3ks := [][]models.PlayerHighlightSchema{
		{
			{RoundNo: 17, HighlightType: models.P3k, PlayerAgainstId: 16459},
			{RoundNo: 17, HighlightType: models.P3k, PlayerAgainstId: 727},
			{RoundNo: 17, HighlightType: models.P3k, PlayerAgainstId: 5654},
		},
		{
			{RoundNo: 20, HighlightType: models.P3k, PlayerAgainstId: 2858},
			{RoundNo: 20, HighlightType: models.P3k, PlayerAgainstId: 727},
			{RoundNo: 20, HighlightType: models.P3k, PlayerAgainstId: 5654},
		},
	}
	p4ks := [][]models.PlayerHighlightSchema{
		{
			{RoundNo: 13, HighlightType: models.P4k, PlayerAgainstId: 727},
			{RoundNo: 13, HighlightType: models.P4k, PlayerAgainstId: 2858},
			{RoundNo: 13, HighlightType: models.P4k, PlayerAgainstId: 5654},
			{RoundNo: 13, HighlightType: models.P4k, PlayerAgainstId: 16459},
		},
	}
	p5ks := [][]models.PlayerHighlightSchema{}
	p1v1s := [][]models.PlayerHighlightSchema{
		{
			{RoundNo: 13, HighlightType: models.P1v1, PlayerAgainstId: 16459},
		},
	}
	p1v2s := [][]models.PlayerHighlightSchema{}
	p1v3s := [][]models.PlayerHighlightSchema{}
	p1v4s := [][]models.PlayerHighlightSchema{}
	p1v5s := [][]models.PlayerHighlightSchema{}

	otherTeamHashmap := map[string]int{
		"Avez":     5654,
		"Elite":    16459,
		"marteen":  10307,
		"SUYGETSU": 2858,
		"saadhak":  727,
	}

	res, err := http.Get(
		"https://www.vlr.gg/510149/fnatic-vs-karmine-corp-esports-world-cup-2025-qf/?tab=performance",
	)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
	sc.Handle(regexp.MustCompile("highlights"), Handler)

	p2kNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(3) > div > div > div",
	)
	p3kNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(4) > div > div > div",
	)
	p4kNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(5) > div > div > div",
	)
	p5kNode := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(6) > div > div > div",
	)
	p1v1Node := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(7) > div > div > div",
	)
	p1v2Node := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(8) > div > div > div",
	)
	p1v3Node := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(9) > div > div > div",
	)
	p1v4Node := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(10) > div > div > div",
	)
	p1v5Node := doc.Find(
		"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div:nth-child(4) > div:nth-child(2) > table > tbody > tr:nth-child(6) > td:nth-child(11) > div > div > div",
	)

	ss := func(testHighlights [][]models.PlayerHighlightSchema, node *goquery.Selection, highlightType models.HighlightType) {
		for i, testHighlight := range testHighlights {
			highlightNode := node.Children().Eq(i)

			data := Data{
				MatchId:          testHighlight[0].MatchId,
				MapId:            testHighlight[0].MapId,
				TeamId:           testHighlight[0].TeamId,
				PlayerId:         testHighlight[0].PlayerId,
				HighlightType:    highlightType,
				OtherTeamHashMap: otherTeamHashmap,
			}

			ctx := context.WithValue(context.Background(), "data", &data)
			ctx2 := context.WithValue(ctx, "tx", tx)

			if err := sc.Pipe("highlights", ctx2, highlightNode); err != nil {
				t.Fatal(err)
			}

			if err := compareHighlights(testHighlight, data.HighlightLog); err != nil {
				t.Error(err)
			}
		}
	}

	ss(p2ks, p2kNode, models.P2k)
	ss(p3ks, p3kNode, models.P3k)
	ss(p4ks, p4kNode, models.P4k)
	ss(p5ks, p5kNode, models.P5k)
	ss(p1v1s, p1v1Node, models.P1v1)
	ss(p1v2s, p1v2Node, models.P1v2)
	ss(p1v3s, p1v3Node, models.P1v3)
	ss(p1v4s, p1v4Node, models.P1v4)
	ss(p1v5s, p1v5Node, models.P1v5)
}
