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
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/playerhighlights"
	"gorm.io/gorm"
)

const (
	playerHighlightRowsSelector = "div:nth-child(2) > table > tbody > tr:not(tr:first-child)"
)

func scrapePlayersHighlights(
	tx *gorm.DB,
	sc *piper.Scraper,
	matchMapSchema models.MatchMapSchema,
	mapPerformanceNode *goquery.Selection,
	t1Hashmap,
	t2Hashmap map[string]int,
) error {
	playerHighlightRows := mapPerformanceNode.Find(playerHighlightRowsSelector)

	if err := tx.Transaction(func(ptx *gorm.DB) error {
		for i := range playerHighlightRows.Length() {
			playerHighlightRow := playerHighlightRows.Eq(i)

			pNameNode := playerHighlightRow.Find(`td:nth-child(1)`)
			p2kNode := playerHighlightRow.Find(`td:nth-child(3) > div > div > div`)
			p3kNode := playerHighlightRow.Find(`td:nth-child(4) > div > div > div`)
			p4kNode := playerHighlightRow.Find(`td:nth-child(5) > div > div > div`)
			p5kNode := playerHighlightRow.Find(`td:nth-child(6) > div > div > div`)
			p1v1Node := playerHighlightRow.Find(`td:nth-child(7) > div > div > div`)
			p1v2Node := playerHighlightRow.Find(`td:nth-child(8) > div > div > div`)
			p1v3Node := playerHighlightRow.Find(`td:nth-child(9) > div > div > div`)
			p1v4Node := playerHighlightRow.Find(`td:nth-child(10) > div > div > div`)
			p1v5Node := playerHighlightRow.Find(`td:nth-child(11) > div > div > div`)

			var teamId int
			teamHashmap := map[string]int{}
			otherTeamHashmap := map[string]int{}

			switch {
			case 0 <= i && i <= 4:
				teamId = matchMapSchema.Team1Id
				teamHashmap = t1Hashmap
				otherTeamHashmap = t2Hashmap
			case 5 <= i && i <= 9:
				teamId = matchMapSchema.Team2Id
				teamHashmap = t2Hashmap
				otherTeamHashmap = t1Hashmap
			}

			playerName := strings.TrimSpace(pNameNode.Find(playerNameSelector).Clone().Children().Remove().End().Text())

			playerId, ok := teamHashmap[playerName]
			if !ok {
				return fmt.Errorf("Player name '%s' doesn't exists in hashmap", playerName)
			}

			errChan := make(chan error)
			doneChan := make(chan bool)

			go func() {
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p2kNode, models.P2k, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p3kNode, models.P3k, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p4kNode, models.P4k, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p5kNode, models.P5k, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p1v1Node, models.P1v1, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p1v2Node, models.P1v2, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p1v3Node, models.P1v3, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p1v4Node, models.P1v4, teamId, playerId, otherTeamHashmap)
				errChan <- scrapeHighlight(tx, sc, matchMapSchema, p1v5Node, models.P1v5, teamId, playerId, otherTeamHashmap)

				doneChan <- true
			}()

			for {
				select {
				case <-doneChan:
					goto CONT
				case err := <-errChan:
					if err != nil {
						return err
					}
				}
			}

		CONT:
			continue
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func scrapeHighlight(
	tx *gorm.DB,
	sc *piper.Scraper,
	matchMapSchema models.MatchMapSchema,
	highlightNode *goquery.Selection,
	highlightType models.HighlightType,
	teamId int,
	playerId int,
	otherTeamHashmap map[string]int,
) error {
	for i := range highlightNode.Children().Length() {
		data := playerhighlights.Data{
			MatchId:          matchMapSchema.MatchId,
			MapId:            matchMapSchema.MapId,
			TeamId:           teamId,
			PlayerId:         playerId,
			HighlightType:    highlightType,
			OtherTeamHashMap: otherTeamHashmap,
		}

		ctx := context.WithValue(context.WithValue(context.Background(), "data", &data), "tx", tx)

		if err := sc.Pipe("highlights", ctx, highlightNode.Children().Eq(i)); err != nil {
			fmt.Printf("player id: %d\n", playerId)
			fmt.Printf("other team hash map: %v\n", otherTeamHashmap)
			return err
		}

		if len(data.HighlightLog) == 0 {
			continue
		}

		prettyPrintHighlight(data.HighlightLog)
	}

	return nil
}

func prettyPrintHighlight(highlightLog []models.PlayerHighlightSchema) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	var row table.Row

	row = append(row, highlightLog[0].RoundNo)
	row = append(row, highlightLog[0].HighlightType)
	row = append(row, highlightLog[0].PlayerId)
	var playersAgainst []string

	for _, fight := range highlightLog {
		playersAgainst = append(playersAgainst, fmt.Sprintf("%d", fight.PlayerAgainstId))
	}

	row = append(row, strings.Join(playersAgainst, ","))

	t.AppendHeader(table.Row{"ROUND", "TYPE", "PLAYER ID", "PLAYERS AGAINST ID"})
	t.AppendRow(row)

	t.Render()
}
