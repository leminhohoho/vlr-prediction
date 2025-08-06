package matchmaps

import (
	"context"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/jedib0t/go-pretty/table"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/playerstats"
	"gorm.io/gorm"
)

const (
	t1PlayersStatsSelector = `div:nth-child(4) > div:nth-child(1) > table > tbody > tr`
	t2PlayersStatsSelector = `div:nth-child(4) > div:nth-child(2) > table > tbody > tr`
)

func scrapePlayersStats(
	tx *gorm.DB,
	sc *piper.Scraper,
	matchMapSchema models.MatchMapSchema,
	mapOverviewNode *goquery.Selection,
) (map[string]int, map[string]int, error) {
	totalOTRounds := matchMapSchema.Team1OTScore + matchMapSchema.Team1OTScore
	t1DefRounds := matchMapSchema.Team1DefScore + matchMapSchema.Team2AtkScore + totalOTRounds/2
	t1AtkRounds := matchMapSchema.Team1AtkScore + matchMapSchema.Team2DefScore + totalOTRounds/2
	t2DefRounds := matchMapSchema.Team2DefScore + matchMapSchema.Team1AtkScore + totalOTRounds/2
	t2AtkRounds := matchMapSchema.Team2AtkScore + matchMapSchema.Team1DefScore + totalOTRounds/2

	t1PlayerStatsNodes := mapOverviewNode.Find(t1PlayersStatsSelector)
	t2PlayerStatsNodes := mapOverviewNode.Find(t2PlayersStatsSelector)

	t1Hashmap := map[string]int{}
	t2Hashmap := map[string]int{}

	pStatsTable := table.NewWriter()
	pStatsTable.SetOutputMirror(os.Stdout)
	pStatsTable.AppendHeader(table.Row{"Rating", "Acs", "K", "D", "A", "KAST", "ADR", "HS", "FK", "FD"})

	for i := range t1PlayerStatsNodes.Length() {
		data := playerstats.Data{
			DefStat: models.PlayerOverviewStatSchema{
				MatchId: matchMapSchema.MatchId,
				MapId:   matchMapSchema.MapId,
				TeamId:  matchMapSchema.Team1Id,
				Side:    models.Def,
			},
			AtkStat: models.PlayerOverviewStatSchema{
				MatchId: matchMapSchema.MatchId,
				MapId:   matchMapSchema.MapId,
				TeamId:  matchMapSchema.Team1Id,
				Side:    models.Atk,
			},
			BothSideStat: models.PlayerOverviewStatSchema{
				Side: models.Side(""),
			},
			TeamDefRounds: t1DefRounds,
			TeamAtkRounds: t1AtkRounds,
		}

		ctx := context.WithValue(context.WithValue(context.Background(), "data", &data), "tx", tx)

		if err := sc.Pipe("playerStats", ctx, t1PlayerStatsNodes.Eq(i)); err != nil {
			return nil, nil, err
		}

		t1Hashmap[data.PlayerName] = data.DefStat.PlayerId
		pStatsTable.AppendRow(table.Row{
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Rating, *data.AtkStat.Rating),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Acs, *data.AtkStat.Acs),
			fmt.Sprintf("%d-%d", *data.DefStat.Kills, *data.AtkStat.Kills),
			fmt.Sprintf("%d-%d", *data.DefStat.Deaths, *data.AtkStat.Deaths),
			fmt.Sprintf("%d-%d", *data.DefStat.Assists, *data.AtkStat.Assists),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Kast, *data.AtkStat.Kast),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Adr, *data.AtkStat.Adr),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Hs, *data.AtkStat.Hs),
			fmt.Sprintf("%d-%d", *data.DefStat.FirstKills, *data.AtkStat.FirstKills),
			fmt.Sprintf("%d-%d", *data.DefStat.FirstDeaths, *data.AtkStat.FirstDeaths),
		})
	}

	pStatsTable.AppendRow(table.Row{})

	for i := range t2PlayerStatsNodes.Length() {
		data := playerstats.Data{
			DefStat: models.PlayerOverviewStatSchema{
				MatchId: matchMapSchema.MatchId,
				MapId:   matchMapSchema.MapId,
				TeamId:  matchMapSchema.Team2Id,
				Side:    models.Def,
			},
			AtkStat: models.PlayerOverviewStatSchema{
				MatchId: matchMapSchema.MatchId,
				MapId:   matchMapSchema.MapId,
				TeamId:  matchMapSchema.Team2Id,
				Side:    models.Atk,
			},
			BothSideStat: models.PlayerOverviewStatSchema{
				Side: models.Side(""),
			},
			TeamDefRounds: t2DefRounds,
			TeamAtkRounds: t2AtkRounds,
		}

		ctx := context.WithValue(context.WithValue(context.Background(), "data", &data), "tx", tx)

		if err := sc.Pipe("playerStats", ctx, t2PlayerStatsNodes.Eq(i)); err != nil {
			return nil, nil, err
		}

		t2Hashmap[data.PlayerName] = data.DefStat.PlayerId
		pStatsTable.AppendRow(table.Row{
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Rating, *data.AtkStat.Rating),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Acs, *data.AtkStat.Acs),
			fmt.Sprintf("%d-%d", *data.DefStat.Kills, *data.AtkStat.Kills),
			fmt.Sprintf("%d-%d", *data.DefStat.Deaths, *data.AtkStat.Deaths),
			fmt.Sprintf("%d-%d", *data.DefStat.Assists, *data.AtkStat.Assists),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Kast, *data.AtkStat.Kast),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Adr, *data.AtkStat.Adr),
			fmt.Sprintf("%.2f-%.2f", *data.DefStat.Hs, *data.AtkStat.Hs),
			fmt.Sprintf("%d-%d", *data.DefStat.FirstKills, *data.AtkStat.FirstKills),
			fmt.Sprintf("%d-%d", *data.DefStat.FirstDeaths, *data.AtkStat.FirstDeaths),
		})
	}

	pStatsTable.Render()

	return t1Hashmap, t2Hashmap, nil
}
