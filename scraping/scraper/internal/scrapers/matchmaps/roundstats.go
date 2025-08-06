package matchmaps

import (
	"context"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/table"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func scrapeRoundsStats(
	tx *gorm.DB,
	sc *piper.Scraper,
	matchMapSchema models.MatchMapSchema,
	mapOverviewNode *goquery.Selection,
	mapEconomyNode *goquery.Selection,
) error {
	var roundNos table.Row
	var team1Banks table.Row
	var team1BuyTypes table.Row
	var team2Banks table.Row
	var team2BuyTypes table.Row
	var teamsDef table.Row
	var teamsWon table.Row
	var wonMethods table.Row

	roundsTable := table.NewWriter()
	roundsTable.SetOutputMirror(os.Stdout)

	roundsOverviewNodes := mapOverviewNode.Find(roundOverviewSelector)
	roundsEconomyNodes := mapEconomyNode.Find(roundEconomySelector)

	if err := tx.Transaction(func(rtx *gorm.DB) error {

		for i := range roundsOverviewNodes.Length() {
			roundOverviewNode := roundsOverviewNodes.Eq(i)
			roundEconomyNode := roundsEconomyNodes.Eq(i)

			combined := roundOverviewNode.Clone().AddSelection(roundEconomyNode.Clone())

			roundStat := models.RoundStatSchema{
				MatchId: matchMapSchema.MatchId,
				MapId:   matchMapSchema.MapId,
				Team1Id: matchMapSchema.Team1Id,
				Team2Id: matchMapSchema.Team2Id,
			}

			roundCtx := context.WithValue(context.WithValue(context.Background(), "roundStat", &roundStat), "tx", rtx)

			if err := sc.Pipe("roundStat", roundCtx, combined); err != nil {
				return err
			}

			roundNos = append(roundNos, roundStat.RoundNo)
			team1Banks = append(team1Banks, fmt.Sprintf("$%d", roundStat.Team1Bank))
			team1BuyTypes = append(team1BuyTypes, shortenBuyType(roundStat.Team1BuyType))
			team2Banks = append(team2Banks, fmt.Sprintf("$%d", roundStat.Team2Bank))
			team2BuyTypes = append(team2BuyTypes, shortenBuyType(roundStat.Team2BuyType))
			teamsDef = append(teamsDef, roundStat.TeamDef)
			teamsWon = append(teamsWon, roundStat.TeamWon)
			wonMethods = append(wonMethods, shortenWonMethod(roundStat.WonMethod))
		}

		return nil
	}); err != nil {
		logrus.Errorf("Error extracting round stats: %s, rounds stats of this map won't be uploaded", err.Error())
		return nil
	}

	for i := 0; i < roundsOverviewNodes.Length(); i += 24 {
		start := i
		var end int
		if roundsOverviewNodes.Length()-i >= 24 {
			end = i + 24
		} else {
			end = i + roundsOverviewNodes.Length()%24
		}

		roundsTable.AppendHeader(append(table.Row{"ROUNDS"}, roundNos[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"T1 BANKS"}, team1Banks[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"T1 BUY"}, team1BuyTypes[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"T2 BANKS"}, team2Banks[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"T2 BUY"}, team2BuyTypes[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"DEF"}, teamsDef[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"WON"}, teamsWon[start:end]...))
		roundsTable.AppendHeader(append(table.Row{"METHOD"}, wonMethods[start:end]...))

		roundsTable.Render()

		roundsTable = table.NewWriter()
		roundsTable.SetOutputMirror(os.Stdout)
	}

	return nil
}

func shortenBuyType(b models.BuyType) string {
	switch b {
	case models.Pistol:
		return "P"
	case models.Eco:
		return "E"
	case models.SemiEco:
		return "SE"
	case models.SemiBuy:
		return "SB"
	case models.FullBuy:
		return "FB"
	}

	return ""
}

func shortenWonMethod(w models.WonMethod) string {
	switch w {
	case models.Defuse:
		return "DEF"
	case models.Eliminate:
		return "ELIM"
	case models.OutOfTime:
		return "OOT"
	case models.SpikeExplode:
		return "SPIKE"
	}

	return ""
}
