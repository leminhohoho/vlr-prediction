package banpicklog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Action string

const (
	Ban    Action = "ban"
	Pick   Action = "pick"
	Remain Action = "remains"
)

type Turn struct {
	MatchId int
	TeamId  *int
	MapId   int
	Action  Action
}

type Data struct {
	MatchId        int
	Team1Id        int
	Team2Id        int
	Team1Shorthand string
	Team2Shorthand string
	Turns          []Turn
}

type BanPickLogScraper struct {
	Data        Data
	BanPickNote string
	Conn        *sql.DB
	Tx          *gorm.Tx
}

func NewBanPickLogScraper(
	conn *sql.DB,
	tx *gorm.Tx,
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
		Conn:        conn,
	}
}

func (b *BanPickLogScraper) getMapId(mapName string) (int, error) {
	var mapId int

	row := b.Conn.QueryRow("SELECT id FROM maps WHERE name = ?", mapName)
	if err := row.Scan(&mapId); err != nil {
		return -1, err
	}

	return mapId, nil
}

func (b *BanPickLogScraper) parseToTurn(teamShorthand, action, mapName string) (Turn, error) {
	var teamId int

	if teamShorthand == b.Data.Team1Shorthand {
		teamId = b.Data.Team1Id
	} else if teamShorthand == b.Data.Team2Shorthand {
		teamId = b.Data.Team2Id
	} else {
		return Turn{}, fmt.Errorf("Unrecognizable team shorthand: %s", teamShorthand)
	}

	mapId, err := b.getMapId(mapName)
	if err != nil {
		return Turn{}, err
	}

	if Action(action) == Ban {
		return Turn{
			MatchId: b.Data.MatchId,
			TeamId:  &teamId,
			MapId:   mapId,
			Action:  Ban,
		}, nil
	} else if Action(action) == Pick {
		return Turn{
			MatchId: b.Data.MatchId,
			TeamId:  &teamId,
			MapId:   mapId,
			Action:  Pick,
		}, nil
	} else {
		return Turn{}, fmt.Errorf("Unable to recognize action %s", action)
	}
}

func (b *BanPickLogScraper) parseToFinalTurn(mapName, action string) (Turn, error) {
	mapId, err := b.getMapId(mapName)
	if err != nil {
		return Turn{}, err
	}

	if Action(action) != Remain {
		return Turn{}, fmt.Errorf("Unable to recognize action %s", action)
	}

	return Turn{
		MatchId: b.Data.MatchId,
		TeamId:  nil,
		MapId:   mapId,
		Action:  Remain,
	}, nil
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
	for _, turnStr := range turnStrs {
		words := strings.Fields(strings.TrimSpace(turnStr))
		if len(words) == 3 {
			turn, err := b.parseToTurn(words[0], words[1], words[2])
			if err != nil {
				return err
			}
			b.Data.Turns = append(b.Data.Turns, turn)
		} else if len(words) == 2 {
			turn, err := b.parseToFinalTurn(words[0], words[1])
			if err != nil {
				return err
			}

			b.Data.Turns = append(b.Data.Turns, turn)
		} else {
			return fmt.Errorf("Unable to determine the turn from this string: %s", turnStr)
		}
	}

	return nil
}
