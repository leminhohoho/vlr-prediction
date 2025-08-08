package playerhighlights

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	PlayerNameSelector = "div:not(:first-child)"
)

type Data struct {
	MatchId          int
	MapId            int
	RoundNo          int `selector:"div:nth-child(1) > span"`
	TeamId           int
	PlayerId         int
	HighlightType    models.HighlightType
	HighlightLog     []models.PlayerHighlightSchema
	OtherTeamHashMap map[string]int
}

type PlayerHighlightScraper struct {
	Data                Data
	PlayerHighlightNode *goquery.Selection
	Tx                  *gorm.DB
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	data, ok := ctx.Value("data").(*Data)
	if !ok {
		return fmt.Errorf("Unable to find data for player highlight")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find gorm transaction")
	}

	logrus.Debug("Getting player highlight round no")
	if err := htmlx.ParseFromSelection(data, selection, htmlx.SetNoPassThroughStruct(true)); err != nil {
		return err
	}

	logrus.Debug("Getting players against ids")
	playersNames := selection.Find(PlayerNameSelector)

	for i := range playersNames.Length() {
		playerNameNode := playersNames.Eq(i)
		playerAgainstName := strings.TrimSpace(playerNameNode.Children().Remove().End().Text())
		if playerAgainstName == "" {
			return fmt.Errorf("Player number %d int the highlight log is empty", i)

		}

		playerAgainstId, ok := data.OtherTeamHashMap[playerAgainstName]
		if !ok {
			logrus.Warnf("Player %s is not in the other team, set id to 0", playerAgainstName)
			playerAgainstId = 0
		}

		data.HighlightLog = append(data.HighlightLog, models.PlayerHighlightSchema{
			MatchId:         data.MatchId,
			MapId:           data.MapId,
			RoundNo:         data.RoundNo,
			TeamId:          data.TeamId,
			PlayerId:        data.PlayerId,
			HighlightType:   data.HighlightType,
			PlayerAgainstId: playerAgainstId,
		})
	}

	logrus.Debug("Saving highlights to db")
	for _, highlight := range data.HighlightLog {
		if err := tx.Table("player_highlights").Create(&highlight).Error; err != nil {
			return err
		}
	}

	return nil
}
