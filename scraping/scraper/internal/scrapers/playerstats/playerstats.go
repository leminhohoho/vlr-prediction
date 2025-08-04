package playerstats

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Data struct {
	DefStat       models.PlayerOverviewStatSchema
	AtkStat       models.PlayerOverviewStatSchema
	BothSideStat  models.PlayerOverviewStatSchema
	TeamDefRounds int
	TeamAtkRounds int
	PlayerName    string `selector:"td.mod-player > div > a > div:nth-child(1)"`
}

func agentParser(tx *gorm.DB) htmlx.Parser {
	return func(rawVal string) (any, error) {
		agentName := strings.TrimSpace(rawVal)
		var agent models.AgentSchema

		rs := tx.Table("agents").Where("name = ?", agentName).First(&agent)
		if rs.Error != nil {
			return nil, rs.Error
		}

		return agent.Id, nil
	}
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	var err error

	data, ok := ctx.Value("data").(*Data)
	if !ok {
		return fmt.Errorf("Unable to find data for player stat scraper")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find both the transaction")
	}

	defStatNode := selection.Clone()
	defStatNode.Find("span.mod-t, span.mod-both").Remove()

	atkStatNode := selection.Clone()
	atkStatNode.Find("span.mod-ct, span.mod-both").Remove()

	bothSideStatNode := selection.Clone()
	bothSideStatNode.Find("span.mod-ct, span.mod-t").Remove()

	// WARNING: Player name will be extracted from match map scraper, not here anymore

	parsers := map[string]htmlx.Parser{
		"agentParser":    agentParser(tx),
		"playerIdParser": customparsers.IdParser,
	}

	logrus.Debug("Parsing player name")
	if err := htmlx.ParseFromSelection(data, selection, htmlx.SetNoPassThroughStruct(true)); err != nil {
		return err
	}

	logrus.Debug("Parsing player def stats")
	if err := htmlx.ParseFromSelection(&data.DefStat, defStatNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Parsing player atk stats")
	if err := htmlx.ParseFromSelection(&data.AtkStat, atkStatNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Parsing player both side stats")
	if err := htmlx.ParseFromSelection(&data.BothSideStat, bothSideStatNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player rating (if any)")
	if data.DefStat.Rating, data.AtkStat.Rating, err = helpers.FillPlayerPerRoundStat(data.DefStat.Rating, data.AtkStat.Rating, data.BothSideStat.Rating, data.TeamDefRounds, data.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player acs (if any)")
	if data.DefStat.Acs, data.AtkStat.Acs, err = helpers.FillPlayerPerRoundStat(data.DefStat.Acs, data.AtkStat.Acs, data.BothSideStat.Acs, data.TeamDefRounds, data.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player kills (if any)")
	if data.DefStat.Kills, data.AtkStat.Kills, err = helpers.FillPlayerKDA(data.DefStat.Kills, data.AtkStat.Kills, data.BothSideStat.Kills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player deaths (if any)")
	if data.DefStat.Deaths, data.AtkStat.Deaths, err = helpers.FillPlayerKDA(data.DefStat.Deaths, data.AtkStat.Deaths, data.BothSideStat.Deaths); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player assists (if any)")
	if data.DefStat.Assists, data.AtkStat.Assists, err = helpers.FillPlayerKDA(data.DefStat.Assists, data.AtkStat.Assists, data.BothSideStat.Assists); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player kast (if any)")
	if data.DefStat.Kast, data.AtkStat.Kast, err = helpers.FillPlayerPerRoundStat(data.DefStat.Kast, data.AtkStat.Kast, data.BothSideStat.Kast, data.TeamDefRounds, data.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player adr (if any)")
	if data.DefStat.Adr, data.AtkStat.Adr, err = helpers.FillPlayerPerRoundStat(data.DefStat.Adr, data.AtkStat.Adr, data.BothSideStat.Adr, data.TeamDefRounds, data.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player hs (if any)")
	if data.DefStat.Hs, data.AtkStat.Hs, err = helpers.FillPlayerPerKillStat(data.DefStat.Hs, data.AtkStat.Hs, data.BothSideStat.Hs, *data.DefStat.Kills, *data.AtkStat.Kills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player fk (if any)")
	if data.DefStat.FirstKills, data.AtkStat.FirstKills, err = helpers.FillPlayerKDA(data.DefStat.FirstKills, data.AtkStat.FirstKills, data.BothSideStat.FirstKills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player fd (if any)")
	if data.DefStat.FirstDeaths, data.AtkStat.FirstDeaths, err = helpers.FillPlayerKDA(data.DefStat.FirstDeaths, data.AtkStat.FirstDeaths, data.BothSideStat.FirstDeaths); err != nil {
		return err
	}

	if err := helpers.PrettyPrintStruct(*data); err != nil {
		return err
	}

	return nil
}
