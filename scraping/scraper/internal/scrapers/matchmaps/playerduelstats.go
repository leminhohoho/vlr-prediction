package matchmaps

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/table"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/gorm"
)

const (
	duelKillsTableSelector      = `div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-normal > tbody`
	duelFirstKillsTableSelector = `div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-fkfd > tbody`
	duelOpKillsTableSelector    = `div:nth-child(1) > table.wf-table-inset.mod-matrix.mod-op > tbody`
	nodeSelector                = `tr:nth-child(%d) > td:nth-child(%d)`
	playerNameSelector          = `div.team > div`
)

func scrapePlayerDuelStats(
	tx *gorm.DB,
	sc *piper.Scraper,
	matchMapSchema models.MatchMapSchema,
	mapPerformanceNode *goquery.Selection,
	t1Hashmap map[string]int,
	t2Hashmap map[string]int,
) error {
	duelKillsTableNode := mapPerformanceNode.Find(duelKillsTableSelector)
	duelFirstKillsTableNode := mapPerformanceNode.Find(duelFirstKillsTableSelector)
	duelOpKillsTableNode := mapPerformanceNode.Find(duelOpKillsTableSelector)

	duelTable := [5][5]models.PlayerDuelStatSchema{}

	if err := tx.Transaction(func(ptx *gorm.DB) error {
		for i := 2; i <= 6; i++ {
			for j := 2; j <= 6; j++ {
				duelKillsNode := duelKillsTableNode.Find(fmt.Sprintf(nodeSelector, i, j))
				duelFirstKillsNode := duelFirstKillsTableNode.Find(fmt.Sprintf(nodeSelector, i, j))
				duelOpKillsNode := duelOpKillsTableNode.Find(fmt.Sprintf(nodeSelector, i, j))

				t1PlayerNode := duelKillsTableNode.Find(fmt.Sprintf(nodeSelector, i, 1))
				t2PlayerNode := duelKillsTableNode.Find(fmt.Sprintf(nodeSelector, 1, j))

				t1PlayerName := strings.TrimSpace(t1PlayerNode.Find(playerNameSelector).Clone().Children().Remove().End().Text())
				t2PlayerName := strings.TrimSpace(t2PlayerNode.Find(playerNameSelector).Clone().Children().Remove().End().Text())

				t1PlayerId, ok := t1Hashmap[t1PlayerName]
				if !ok {
					return fmt.Errorf("Player name '%s' doesn't exists in t1 hashmap", t1PlayerName)
				}

				t2PlayerId, ok := t2Hashmap[t2PlayerName]
				if !ok {
					return fmt.Errorf("Player name '%s' doesn't exists in t2 hashmap", t2PlayerName)
				}

				combined := duelKillsNode.Clone().AddSelection(duelFirstKillsNode).AddSelection(duelOpKillsNode)

				duelStats := models.PlayerDuelStatSchema{
					MatchId:       matchMapSchema.MatchId,
					MapId:         matchMapSchema.MapId,
					Team1PlayerId: t1PlayerId,
					Team2PlayerId: t2PlayerId,
				}

				ctx := context.WithValue(context.WithValue(context.Background(), "duelStats", &duelStats), "tx", ptx)

				if err := sc.Pipe("duelStats", ctx, combined); err != nil {
					return err
				}

				duelTable[i-2][j-2] = duelStats
			}
		}

		return nil
	}); err != nil {
		return err
	}

	duelKillsTable := table.NewWriter()
	duelFirstKillsTable := table.NewWriter()
	duelOpKillsTable := table.NewWriter()
	duelKillsTable.SetOutputMirror(os.Stdout)
	duelFirstKillsTable.SetOutputMirror(os.Stdout)
	duelOpKillsTable.SetOutputMirror(os.Stdout)

	duelKillsTable.AppendHeader(table.Row{
		"",
		duelTable[0][0].Team2PlayerId,
		duelTable[0][1].Team2PlayerId,
		duelTable[0][2].Team2PlayerId,
		duelTable[0][3].Team2PlayerId,
		duelTable[0][4].Team2PlayerId,
	})
	duelFirstKillsTable.AppendHeader(table.Row{
		"",
		duelTable[0][0].Team2PlayerId,
		duelTable[0][1].Team2PlayerId,
		duelTable[0][2].Team2PlayerId,
		duelTable[0][3].Team2PlayerId,
		duelTable[0][4].Team2PlayerId,
	})
	duelOpKillsTable.AppendHeader(table.Row{
		"",
		duelTable[0][0].Team2PlayerId,
		duelTable[0][1].Team2PlayerId,
		duelTable[0][2].Team2PlayerId,
		duelTable[0][3].Team2PlayerId,
		duelTable[0][4].Team2PlayerId,
	})

	for i := range 5 {
		duelKillsTable.AppendRow(table.Row{
			duelTable[i][0].Team1PlayerId,
			fmt.Sprintf("%d-%d", duelTable[i][0].Team1PlayerKillsVsTeam2Player, duelTable[i][0].Team2PlayerKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][1].Team1PlayerKillsVsTeam2Player, duelTable[i][1].Team2PlayerKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][2].Team1PlayerKillsVsTeam2Player, duelTable[i][2].Team2PlayerKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][3].Team1PlayerKillsVsTeam2Player, duelTable[i][3].Team2PlayerKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][4].Team1PlayerKillsVsTeam2Player, duelTable[i][4].Team2PlayerKillsVsTeam1Player),
		})
		duelFirstKillsTable.AppendRow(table.Row{
			duelTable[i][0].Team1PlayerId,
			fmt.Sprintf("%d-%d", duelTable[i][0].Team1PlayerFirstKillsVsTeam2Player, duelTable[i][0].Team2PlayerFirstKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][1].Team1PlayerFirstKillsVsTeam2Player, duelTable[i][1].Team2PlayerFirstKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][2].Team1PlayerFirstKillsVsTeam2Player, duelTable[i][2].Team2PlayerFirstKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][3].Team1PlayerFirstKillsVsTeam2Player, duelTable[i][3].Team2PlayerFirstKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][4].Team1PlayerFirstKillsVsTeam2Player, duelTable[i][4].Team2PlayerFirstKillsVsTeam1Player),
		})
		duelOpKillsTable.AppendRow(table.Row{
			duelTable[i][0].Team1PlayerId,
			fmt.Sprintf("%d-%d", duelTable[i][0].Team1PlayerOpKillsVsTeam2Player, duelTable[i][0].Team2PlayerOpKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][1].Team1PlayerOpKillsVsTeam2Player, duelTable[i][1].Team2PlayerOpKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][2].Team1PlayerOpKillsVsTeam2Player, duelTable[i][2].Team2PlayerOpKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][3].Team1PlayerOpKillsVsTeam2Player, duelTable[i][3].Team2PlayerOpKillsVsTeam1Player),
			fmt.Sprintf("%d-%d", duelTable[i][4].Team1PlayerOpKillsVsTeam2Player, duelTable[i][4].Team2PlayerOpKillsVsTeam1Player),
		})
	}

	fmt.Println("== DUEL KILLS TABLE ==")
	duelKillsTable.Render()
	fmt.Println("== DUEL FIRST KILLS TABLE ==")
	duelFirstKillsTable.Render()
	fmt.Println("== DUEL OP KILLS TABLE ==")
	duelOpKillsTable.Render()

	return nil
}
