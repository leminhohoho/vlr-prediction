package playerhighlights

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HighlightType string

const (
	P2k  HighlightType = "2k"
	P3k  HighlightType = "3k"
	P4k  HighlightType = "4k"
	P5k  HighlightType = "5k"
	P1v1 HighlightType = "1v1"
	P1v2 HighlightType = "1v2"
	P1v3 HighlightType = "1v3"
	P1v4 HighlightType = "1v4"
	P1v5 HighlightType = "1v5"

	PlayerNameSelector = "div:not(:first-child)"
)

type PlayerHighlightSchema struct {
	MatchId         int
	MapId           int
	RoundNo         int
	TeamId          int
	PlayerId        int
	HighlightType   HighlightType
	PlayerAgainstId int
}

type Data struct {
	MatchId          int
	MapId            int
	RoundNo          int `selector:"div:nth-child(1) > span"`
	TeamId           int
	PlayerId         int
	HighlightType    HighlightType
	HighlightLog     []PlayerHighlightSchema
	OtherTeamHashMap map[string]int
}

type PlayerHighlightScraper struct {
	Data                Data
	PlayerHighlightNode *goquery.Selection
	Conn                *sql.DB
	Tx                  *gorm.Tx
}

func NewPlayerHighlightScraper(
	conn *sql.DB,
	tx *gorm.Tx,
	playerHighlightNode *goquery.Selection,
	matchId, mapId, teamId, playerId int,
	highlightType HighlightType,
	otherTeamHashMap map[string]int,
) *PlayerHighlightScraper {
	return &PlayerHighlightScraper{
		Data: Data{
			MatchId:          matchId,
			MapId:            mapId,
			TeamId:           teamId,
			PlayerId:         playerId,
			HighlightType:    highlightType,
			OtherTeamHashMap: otherTeamHashMap,
		},
		PlayerHighlightNode: playerHighlightNode,
		Conn:                conn,
		Tx:                  tx,
	}
}

func (p *PlayerHighlightScraper) getPlayersId() error {
	errChan := make(chan error)

	go func() {
		p.PlayerHighlightNode.Find(PlayerNameSelector).Each(func(i int, playerNameNode *goquery.Selection) {
			playerAgainstName := strings.TrimSpace(playerNameNode.Children().Remove().End().Text())
			if playerAgainstName == "" {
				errChan <- fmt.Errorf("Player number %d int the highlight log is empty", i)
				return
			}

			playerAgainstId, ok := p.Data.OtherTeamHashMap[playerAgainstName]
			if !ok {
				errChan <- fmt.Errorf("Player %s is not in the other team", playerAgainstName)
				return
			}

			p.Data.HighlightLog = append(p.Data.HighlightLog, PlayerHighlightSchema{
				MatchId:         p.Data.MatchId,
				MapId:           p.Data.MapId,
				RoundNo:         p.Data.RoundNo,
				TeamId:          p.Data.TeamId,
				PlayerId:        p.Data.PlayerId,
				HighlightType:   p.Data.HighlightType,
				PlayerAgainstId: playerAgainstId,
			})

			errChan <- nil
		})
	}()

	select {
	case err := <-errChan:
		return err
	}
}

func (p *PlayerHighlightScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(p.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (p *PlayerHighlightScraper) Scrape() error {
	logrus.Debug("Getting player highlight round no")
	if err := htmlx.ParseFromSelection(&p.Data, p.PlayerHighlightNode, htmlx.SetNoPassThroughStruct(true)); err != nil {
		return err
	}

	logrus.Debug("Getting players against ids")
	if err := p.getPlayersId(); err != nil {
		return err
	}

	return nil
}
