package playerhighlights

import (
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
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

	for i, p2k := range p2ks {
		highlightNode := p2kNode.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p2k[0].MatchId,
			p2k[0].MapId,
			p2k[0].TeamId,
			p2k[0].PlayerId,
			models.P2k,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := compareHighlights(p2k, s.Data.HighlightLog); err != nil {
			t.Error(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p3k := range p3ks {
		highlightNode := p3kNode.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p3k[0].MatchId,
			p3k[0].MapId,
			p3k[0].TeamId,
			p3k[0].PlayerId,
			models.P3k,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p4k := range p4ks {
		highlightNode := p4kNode.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p4k[0].MatchId,
			p4k[0].MapId,
			p4k[0].TeamId,
			p4k[0].PlayerId,
			models.P4k,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p5k := range p5ks {
		highlightNode := p5kNode.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p5k[0].MatchId,
			p5k[0].MapId,
			p5k[0].TeamId,
			p5k[0].PlayerId,
			models.P5k,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p1v1 := range p1v1s {
		highlightNode := p1v1Node.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p1v1[0].MatchId,
			p1v1[0].MapId,
			p1v1[0].TeamId,
			p1v1[0].PlayerId,
			models.P1v1,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p1v2 := range p1v2s {
		highlightNode := p1v2Node.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p1v2[0].MatchId,
			p1v2[0].MapId,
			p1v2[0].TeamId,
			p1v2[0].PlayerId,
			models.P1v2,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p1v3 := range p1v3s {
		highlightNode := p1v3Node.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p1v3[0].MatchId,
			p1v3[0].MapId,
			p1v3[0].TeamId,
			p1v3[0].PlayerId,
			models.P1v3,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p1v4 := range p1v4s {
		highlightNode := p1v4Node.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p1v4[0].MatchId,
			p1v4[0].MapId,
			p1v4[0].TeamId,
			p1v4[0].PlayerId,
			models.P1v4,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}

	for i, p1v5 := range p1v5s {
		highlightNode := p1v5Node.Children().Eq(i)

		s := NewScraper(
			nil,
			highlightNode,
			p1v5[0].MatchId,
			p1v5[0].MapId,
			p1v5[0].TeamId,
			p1v5[0].PlayerId,
			models.P1v5,
			otherTeamHashmap,
		)

		if err := s.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := s.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
