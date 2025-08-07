package playerduelstats

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var parsers = map[string]htmlx.Parser{
	"duelParser": htmlx.IfNullParser(0, htmlx.IntParser),
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	duelStats, ok := ctx.Value("duelStats").(*models.PlayerDuelStatSchema)
	if !ok {
		return fmt.Errorf("Unable to find player duel stats")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find gorm transaction")
	}

	duelKillsNode := selection.Eq(0)
	duelFirstKillsNode := selection.Eq(1)
	duelOpKillsNode := selection.Eq(2)

	var duelKills models.DuelKills
	var duelFirstKills models.DuelFirstKills
	var duelOpKills models.DuelOpKills

	logrus.Debug("Parsing player duel kills information from html onto match schema")
	if err := htmlx.ParseFromSelection(&duelKills, duelKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Parsing player duel first kills information from html onto match schema")
	if err := htmlx.ParseFromSelection(&duelFirstKills, duelFirstKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	logrus.Debug("Parsing player duel op kills information from html onto match schema")
	if err := htmlx.ParseFromSelection(&duelOpKills, duelOpKillsNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	duelStats.DuelKills = duelKills
	duelStats.DuelFirstKills = duelFirstKills
	duelStats.DuelOpKills = duelOpKills

	logrus.Debug("Saving player duel stats to db")
	if err := tx.Table("players_duel_stats").Create(duelStats).Error; err != nil {
		return err
	}

	return nil
}
