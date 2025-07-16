package playerstats

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initPlayerOverviewStatSchema(
	matchId,
	mapId,
	teamId int,
	side models.Side,
) models.PlayerOverviewStatSchema {
	return models.PlayerOverviewStatSchema{
		MatchId: matchId,
		MapId:   mapId,
		TeamId:  teamId,
		Side:    side,
	}
}

type Data struct {
	DefStat      models.PlayerOverviewStatSchema
	AtkStat      models.PlayerOverviewStatSchema
	BothSideStat models.PlayerOverviewStatSchema
	PlayerName   string `selector:"td.mod-player > div > a > div:nth-child(1)"`
}

type PlayerOverviewStatScraper struct {
	Data                   Data
	PlayerOverviewStatNode *goquery.Selection
	TeamDefRounds          int
	TeamAtkRounds          int
	Tx                     *gorm.DB
}

func NewPlayerOverviewStatScraper(
	tx *gorm.DB,
	playerOverviewStatNode *goquery.Selection,
	matchId int,
	mapId int,
	teamId int,
	// NOTE: def and atk rounds here include OT rounds also
	teamDefRounds int,
	teamAtkRounds int,
) *PlayerOverviewStatScraper {
	return &PlayerOverviewStatScraper{
		Data: Data{
			DefStat:      initPlayerOverviewStatSchema(matchId, mapId, teamId, models.Def),
			AtkStat:      initPlayerOverviewStatSchema(matchId, mapId, teamId, models.Atk),
			BothSideStat: initPlayerOverviewStatSchema(matchId, mapId, teamId, models.Side("")),
		},
		PlayerOverviewStatNode: playerOverviewStatNode,
		TeamDefRounds:          teamDefRounds,
		TeamAtkRounds:          teamAtkRounds,
		Tx:                     tx,
	}
}

func (p *PlayerOverviewStatScraper) agentParser(rawVal string) (any, error) {
	agentName := strings.TrimSpace(rawVal)
	var agent models.AgentSchema

	rs := p.Tx.Table("agents").Where("name = ?", agentName).Scan(&agent)
	if rs.Error != nil {
		return nil, rs.Error
	}

	return agent.Id, nil
}

// NOTE: Split the node in the node for def, atk and both
func (p *PlayerOverviewStatScraper) splitNode() (*goquery.Selection, *goquery.Selection, *goquery.Selection, error) {
	defPlayerOverviewStatNode := p.PlayerOverviewStatNode.Clone()
	defPlayerOverviewStatNode.Find("span.mod-t, span.mod-both").Remove()

	atkPlayerOverviewStatNode := p.PlayerOverviewStatNode.Clone()
	atkPlayerOverviewStatNode.Find("span.mod-ct, span.mod-both").Remove()

	bothSidePlayerOverviewStatNode := p.PlayerOverviewStatNode.Clone()
	bothSidePlayerOverviewStatNode.Find("span.mod-ct, span.mod-t").Remove()

	return defPlayerOverviewStatNode, atkPlayerOverviewStatNode, bothSidePlayerOverviewStatNode, nil
}

func (p *PlayerOverviewStatScraper) fillStats() error {
	var err error

	logrus.Debug("Fill in missing data for player rating (if any)")
	if p.Data.DefStat.Rating, p.Data.AtkStat.Rating, err = helpers.FillPlayerPerRoundStat(p.Data.DefStat.Rating, p.Data.AtkStat.Rating, p.Data.BothSideStat.Rating, p.TeamDefRounds, p.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player acs (if any)")
	if p.Data.DefStat.Acs, p.Data.AtkStat.Acs, err = helpers.FillPlayerPerRoundStat(p.Data.DefStat.Acs, p.Data.AtkStat.Acs, p.Data.BothSideStat.Acs, p.TeamDefRounds, p.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player kills (if any)")
	if p.Data.DefStat.Kills, p.Data.AtkStat.Kills, err = helpers.FillPlayerKDA(p.Data.DefStat.Kills, p.Data.AtkStat.Kills, p.Data.BothSideStat.Kills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player deaths (if any)")
	if p.Data.DefStat.Deaths, p.Data.AtkStat.Deaths, err = helpers.FillPlayerKDA(p.Data.DefStat.Deaths, p.Data.AtkStat.Deaths, p.Data.BothSideStat.Deaths); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player assists (if any)")
	if p.Data.DefStat.Assists, p.Data.AtkStat.Assists, err = helpers.FillPlayerKDA(p.Data.DefStat.Assists, p.Data.AtkStat.Assists, p.Data.BothSideStat.Assists); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player kast (if any)")
	if p.Data.DefStat.Kast, p.Data.AtkStat.Kast, err = helpers.FillPlayerPerRoundStat(p.Data.DefStat.Kast, p.Data.AtkStat.Kast, p.Data.BothSideStat.Kast, p.TeamDefRounds, p.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player adr (if any)")
	if p.Data.DefStat.Adr, p.Data.AtkStat.Adr, err = helpers.FillPlayerPerRoundStat(p.Data.DefStat.Adr, p.Data.AtkStat.Adr, p.Data.BothSideStat.Adr, p.TeamDefRounds, p.TeamAtkRounds); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player hs (if any)")
	if p.Data.DefStat.Hs, p.Data.AtkStat.Hs, err = helpers.FillPlayerPerKillStat(p.Data.DefStat.Hs, p.Data.AtkStat.Hs, p.Data.BothSideStat.Hs, *p.Data.DefStat.Kills, *p.Data.AtkStat.Kills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player fk (if any)")
	if p.Data.DefStat.FirstKills, p.Data.AtkStat.FirstKills, err = helpers.FillPlayerKDA(p.Data.DefStat.FirstKills, p.Data.AtkStat.FirstKills, p.Data.BothSideStat.FirstKills); err != nil {
		return err
	}

	logrus.Debug("Fill in missing data for player fd (if any)")
	if p.Data.DefStat.FirstDeaths, p.Data.AtkStat.FirstDeaths, err = helpers.FillPlayerKDA(p.Data.DefStat.FirstDeaths, p.Data.AtkStat.FirstDeaths, p.Data.BothSideStat.FirstDeaths); err != nil {
		return err
	}

	return nil
}

func (p *PlayerOverviewStatScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(p.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (p *PlayerOverviewStatScraper) Scrape() error {
	defPlayerOverviewStatNode, atkPlayerOverviewStatNode, bothSidePlayerOverviewStatNode, err := p.splitNode()
	if err != nil {
		return err
	}

	logrus.Debug("Parsing player name")
	if err := htmlx.ParseFromSelection(&p.Data, defPlayerOverviewStatNode, htmlx.SetNoPassThroughStruct(true)); err != nil {
		return err
	}

	parsers := map[string]htmlx.Parser{
		"agentParser":    p.agentParser,
		"playerIdParser": customparsers.IdParser,
	}

	logrus.Debug("Parsing player def stats")
	if err := htmlx.ParseFromSelection(
		&p.Data.DefStat,
		defPlayerOverviewStatNode,
		htmlx.SetParsers(parsers),
	); err != nil {
		return err
	}

	logrus.Debug("Parsing player atk stats")
	if err := htmlx.ParseFromSelection(
		&p.Data.AtkStat,
		atkPlayerOverviewStatNode,
		htmlx.SetParsers(parsers),
	); err != nil {
		return err
	}

	logrus.Debug("Parsing player both side stats")
	if err := htmlx.ParseFromSelection(
		&p.Data.BothSideStat,
		bothSidePlayerOverviewStatNode,
		htmlx.SetParsers(parsers),
	); err != nil {
		return err
	}

	logrus.Debug("Fill in any missing player stats")
	if err := p.fillStats(); err != nil {
		return err
	}

	return nil
}
