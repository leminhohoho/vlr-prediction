package roundstats

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

type BuyType string

type WonMethod string

const (
	Pistol  BuyType = "pistol"
	Eco     BuyType = "eco"
	SemiEco BuyType = "semi_eco"
	SemiBuy BuyType = "semi_buy"
	FullBuy BuyType = "full_buy"

	Eliminate    WonMethod = "eliminate"
	SpikeExplode WonMethod = "spike_explode"
	Defuse       WonMethod = "defuse"
	OutOfTime    WonMethod = "out_of_time"
)

type RoundOverviewSchema struct {
	RoundNo   int       `selector:"div.rnd-num"`
	TeamWon   int       `selector:"div.rnd-sq.mod-win:nth-child(2)" source:"attr=class" parser:"teamWonParser"`
	WonMethod WonMethod `selector:"div.rnd-sq.mod-win > img"        source:"attr=src"   parser:"wonMethodParser"`
}

type RoundEconomySchema struct {
	TeamDef      int     `selector:"div.rnd-sq.mod-win"      source:"attr=class" parser:"teamDefParser"`
	Team1BuyType BuyType `selector:"div.rnd-sq:nth-child(3)"                     parser:"buyTypeParser"`
	Team2BuyType BuyType `selector:"div.rnd-sq:nth-child(4)"                     parser:"buyTypeParser"`
	Team1Bank    int     `selector:"div.bank:nth-child(2)"                       parser:"balanceParser"`
	Team2Bank    int     `selector:"div.bank:nth-child(5)"                       parser:"balanceParser"`
}

type RoundStatSchema struct {
	MatchId int
	MapId   int
	Team1Id int
	Team2Id int
	RoundOverviewSchema
	RoundEconomySchema
}

type RoundStatScraper struct {
	Data              RoundStatSchema
	RoundOverviewNode *goquery.Selection
	RoundEconomyNode  *goquery.Selection
	Conn              *sql.DB
	Tx                *gorm.Tx
}

func NewRoundStatScraper(
	conn *sql.DB,
	tx *gorm.Tx,
	roundOverviewNode *goquery.Selection,
	roundEconomyNode *goquery.Selection,
	matchId int,
	mapId int,
	team1Id int,
	team2Id int,
) *RoundStatScraper {
	return &RoundStatScraper{
		Data: RoundStatSchema{
			MatchId: matchId,
			MapId:   mapId,
			Team1Id: team1Id,
			Team2Id: team2Id,
		},
		RoundOverviewNode: roundOverviewNode,
		RoundEconomyNode:  roundEconomyNode,
		Conn:              conn,
		Tx:                tx,
	}
}

func (r *RoundStatScraper) teamWonParser(rawVal string) (any, error) {
	if strings.TrimSpace(rawVal) == "" {
		return r.Data.Team2Id, nil
	}

	return r.Data.Team1Id, nil
}

func wonMethodParser(rawVal string) (any, error) {
	src := strings.TrimSpace(rawVal)
	switch src {
	case "/img/vlr/game/round/elim.webp":
		return Eliminate, nil
	case "/img/vlr/game/round/boom.webp":
		return SpikeExplode, nil
	case "/img/vlr/game/round/defuse.webp":
		return Defuse, nil
	case "/img/vlr/game/round/time.webp":
		return OutOfTime, nil
	default:
		return nil, fmt.Errorf("Unable to specify the won method from this img src: %s", src)
	}
}

// NOTE: This parse must be run after team won has been retrieved (parser for overview need to run first)
func (r *RoundStatScraper) teamDefParser(rawVal string) (any, error) {
	teamWonClasses := strings.TrimSpace(rawVal)
	var teamWon, teamLost int
	if r.Data.Team1Id == r.Data.TeamWon {
		teamWon = r.Data.Team1Id
		teamLost = r.Data.Team2Id
	} else {
		teamWon = r.Data.Team2Id
		teamLost = r.Data.Team1Id
	}

	if strings.Contains(teamWonClasses, "mod-ct") {
		return teamWon, nil
	} else if strings.Contains(teamWonClasses, "mod-t") {
		return teamLost, nil
	} else {
		return nil, fmt.Errorf("Error determine the team def from this class: %s", teamWonClasses)
	}
}

func (r *RoundStatScraper) buyTypeParser(rawVal string) (any, error) {
	buyStr := strings.TrimSpace(rawVal)

	if r.Data.RoundNo == 1 || r.Data.RoundNo == 12 {
		return Pistol, nil
	}

	switch buyStr {
	case "":
		return Eco, nil
	case "$":
		return SemiEco, nil
	case "$$":
		return SemiBuy, nil
	case "$$$":
		return FullBuy, nil
	default:
		return nil, fmt.Errorf("Unable to determinte the buy type from this string: %s", buyStr)
	}
}

func balanceParser(rawVal string) (any, error) {
	balanceVal, err := htmlx.FloatParser(rawVal)
	if err != nil {
		return nil, err
	}

	floatBalance, _ := balanceVal.(float64)
	balance := int(floatBalance * 1000)

	return balance, nil
}

func (r *RoundStatScraper) scrapeOverviewInfo() error {
	logrus.Debug("Scraping round overview info")

	var roundOverviewSchema RoundOverviewSchema

	parsers := map[string]htmlx.Parser{
		"teamWonParser":   r.teamWonParser,
		"wonMethodParser": wonMethodParser,
	}

	if err := htmlx.ParseFromSelection(&roundOverviewSchema, r.RoundOverviewNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	r.Data.RoundOverviewSchema = roundOverviewSchema

	return nil
}

func (r *RoundStatScraper) scrapeEconomyInfo() error {
	logrus.Debug("Scraping round economy info")

	var roundEconomySchema RoundEconomySchema

	parsers := map[string]htmlx.Parser{
		"teamDefParser": r.teamDefParser,
		"buyTypeParser": r.buyTypeParser,
		"balanceParser": balanceParser,
	}

	if err := htmlx.ParseFromSelection(&roundEconomySchema, r.RoundEconomyNode, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	r.Data.RoundEconomySchema = roundEconomySchema

	return nil
}

func (r *RoundStatScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(r.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (r *RoundStatScraper) Scrape() error {
	if err := r.scrapeOverviewInfo(); err != nil {
		return err
	}

	if err := r.scrapeEconomyInfo(); err != nil {
		return err
	}

	return nil
}
