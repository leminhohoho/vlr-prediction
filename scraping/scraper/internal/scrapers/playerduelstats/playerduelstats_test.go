package playerduelstats

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
)

func newTestPlayerDuelStat(
	p1KillsVsP2, p1FKillsVsP2, p1OpKillsVsP2, p2KillsVsP1, p2FKillsVsP1, p2OpKillsVsP1 int,
) models.PlayerDuelStatSchema {
	return models.PlayerDuelStatSchema{
		DuelKills: models.DuelKills{
			Team1PlayerKillsVsTeam2Player: p1KillsVsP2,
			Team2PlayerKillsVsTeam1Player: p2KillsVsP1,
		},
		DuelFirstKills: models.DuelFirstKills{
			Team1PlayerFirstKillsVsTeam2Player: p1FKillsVsP2,
			Team2PlayerFirstKillsVsTeam1Player: p2FKillsVsP1,
		},
		DuelOpKills: models.DuelOpKills{
			Team1PlayerOpKillsVsTeam2Player: p1OpKillsVsP2,
			Team2PlayerOpKillsVsTeam1Player: p2OpKillsVsP1,
		},
	}
}

func TestPlayerDuelStat(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	testDuelStats := []models.PlayerDuelStatSchema{
		newTestPlayerDuelStat(2, 0, 0, 2, 1, 0),
		newTestPlayerDuelStat(4, 1, 0, 3, 1, 0),
		newTestPlayerDuelStat(1, 1, 0, 1, 0, 0),
		newTestPlayerDuelStat(4, 1, 0, 2, 0, 0),
		newTestPlayerDuelStat(1, 0, 0, 2, 0, 0),
	}

	res, err := http.Get(
		"https://www.vlr.gg/510154/gen-g-vs-team-heretics-esports-world-cup-2025-sf/?tab=performance",
	)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	for i, testDuelStat := range testDuelStats {
		duelKillsNode := doc.Find(
			fmt.Sprintf(
				`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='225042'] > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-normal > tbody > tr:nth-child(2) > td:nth-child(%d)`,
				i+2,
			),
		)

		duelFirstKillsNode := doc.Find(
			fmt.Sprintf(
				`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='225042'] > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-fkfd> tbody > tr:nth-child(2) > td:nth-child(%d)`,
				i+2,
			),
		)

		duelOpKillsNode := doc.Find(
			fmt.Sprintf(
				`#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div.vm-stats-game[data-game-id='225042'] > div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-op> tbody > tr:nth-child(2) > td:nth-child(%d)`,
				i+2,
			),
		)

		d := NewPlayerDuelStatScraper(
			nil,
			duelKillsNode,
			duelFirstKillsNode,
			duelOpKillsNode,
			testDuelStat.MatchId,
			testDuelStat.MapId,
			testDuelStat.Team1PlayerId,
			testDuelStat.Team2PlayerId,
		)

		if err := d.Scrape(); err != nil {
			t.Fatal(err)
		}

		if err := helpers.CompareStructs(testDuelStat, d.Data); err != nil {
			t.Error(err)
		}

		if err := d.PrettyPrint(); err != nil {
			t.Fatal(err)
		}
	}
}
