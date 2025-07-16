package playerduelstats

import (
	"encoding/json"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlayerDuelStatScraper struct {
	Data                     models.PlayerDuelStatSchema
	PlayerDuelKillsNode      *goquery.Selection
	PlayerDuelFirstKillsNode *goquery.Selection
	PlayerDuelOpKillsNode    *goquery.Selection
	Tx                       *gorm.Tx
}

func NewPlayerDuelStatScraper(
	tx *gorm.Tx,
	playerDuelKillsNode *goquery.Selection,
	playerDuelFirstKillsNode *goquery.Selection,
	playerDuelOpKillsNode *goquery.Selection,
	matchId, mapId, team1PlayerId, team2PlayerId int,
) *PlayerDuelStatScraper {
	return &PlayerDuelStatScraper{
		Data: models.PlayerDuelStatSchema{
			MatchId:       matchId,
			MapId:         mapId,
			Team1PlayerId: team1PlayerId,
			Team2PlayerId: team2PlayerId,
		},
		PlayerDuelKillsNode:      playerDuelKillsNode,
		PlayerDuelFirstKillsNode: playerDuelFirstKillsNode,
		PlayerDuelOpKillsNode:    playerDuelOpKillsNode,
		Tx:                       tx,
	}
}

func (p *PlayerDuelStatScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(p.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (p *PlayerDuelStatScraper) Scrape() error {
	parsers := map[string]htmlx.Parser{
		"duelParser": htmlx.IfNullParser(0, htmlx.IntParser),
	}

	var duelKills models.DuelKills
	var duelFirstKills models.DuelFirstKills
	var duelOpKills models.DuelOpKills

	logrus.Debug("Getting player duel kills")
	if err := htmlx.ParseFromSelection(&duelKills, p.PlayerDuelKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Getting player duel first kills")
	if err := htmlx.ParseFromSelection(&duelFirstKills, p.PlayerDuelFirstKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Getting player duel op kills")
	if err := htmlx.ParseFromSelection(&duelOpKills, p.PlayerDuelOpKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	p.Data.DuelKills = duelKills
	p.Data.DuelFirstKills = duelFirstKills
	p.Data.DuelOpKills = duelOpKills

	return nil
}
