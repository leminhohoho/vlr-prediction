package banpicklog

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Data struct {
	MatchId        int
	Team1Id        int
	Team2Id        int
	Team1Shorthand string
	Team2Shorthand string
	Turns          []models.BanPickLogSchema
}

type BanPickLogScraper struct {
	Data        Data
	BanPickNote string
	Tx          *gorm.DB
}

func NewBanPickLogScraper(
	tx *gorm.DB,
	matchId, team1Id, team2Id int,
	team1Shorthand, team2Shorthand, banPickNote string,
) *BanPickLogScraper {
	return &BanPickLogScraper{
		Data: Data{
			MatchId:        matchId,
			Team1Id:        team1Id,
			Team2Id:        team2Id,
			Team1Shorthand: team1Shorthand,
			Team2Shorthand: team2Shorthand,
		},
		BanPickNote: banPickNote,
		Tx:          tx,
	}
}

func (b *BanPickLogScraper) getMapId(mapName string) (int, error) {
	var vlrMap models.MapSchema

	rs := b.Tx.Table("maps").Where("name = ?", mapName).First(&vlrMap)
	if rs.Error != nil {
		return -1, rs.Error
	}

	return vlrMap.Id, nil
}

func (b *BanPickLogScraper) parseToTurn(
	teamShorthand, action, mapName string,
) (banPickLog models.BanPickLogSchema, err error) {
	banPickLog.MatchId = b.Data.MatchId

	switch teamShorthand {
	case b.Data.Team1Shorthand:
		banPickLog.TeamId = &b.Data.Team1Id
	case b.Data.Team2Shorthand:
		banPickLog.TeamId = &b.Data.Team2Id
	default:
		err = fmt.Errorf("Unrecognizable team shorthand: %s", teamShorthand)
		return
	}

	banPickLog.MapId, err = b.getMapId(mapName)
	if err != nil {
		return
	}

	switch models.VetoAction(action) {
	case models.BanMap:
		banPickLog.Action = models.BanMap
	case models.PickMap:
		banPickLog.Action = models.PickMap
	default:
		err = fmt.Errorf("Unable to recognize action %s", action)
		return
	}

	return
}

func (b *BanPickLogScraper) parseToFinalTurn(
	mapName, action string,
) (banPickLog models.BanPickLogSchema, err error) {
	banPickLog.MapId, err = b.getMapId(mapName)
	if err != nil {
		return
	}

	if models.VetoAction(action) != models.RemainMap {
		err = fmt.Errorf("Unable to recognize action %s", action)
		return
	}

	banPickLog.MatchId = b.Data.MatchId
	banPickLog.Action = models.RemainMap

	return
}

func (b *BanPickLogScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(b.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (b *BanPickLogScraper) Scrape() error {
	turnStrs := strings.Split(b.BanPickNote, ";")

	logrus.Debug("Getting ban pick log info")
	for i, turnStr := range turnStrs {
		words := strings.Fields(strings.TrimSpace(turnStr))
		switch len(words) {
		case 3:
			if turn, err := b.parseToTurn(words[0], words[1], words[2]); err != nil {
				return err
			} else {
				turn.Order = i + 1
				b.Data.Turns = append(b.Data.Turns, turn)
			}
		case 2:
			if turn, err := b.parseToFinalTurn(words[0], words[1]); err != nil {
				return err
			} else {
				turn.Order = i + 1
				b.Data.Turns = append(b.Data.Turns, turn)
			}
		default:
			return fmt.Errorf("Unable to determine the turn from this string: %s", turnStr)
		}
	}

	return nil
}
