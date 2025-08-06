package roundstats

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

func teamWonParser(t1Id, t2Id int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		if strings.TrimSpace(rawVal) == "" {
			return t2Id, nil
		}

		return t1Id, nil
	}
}

func wonMethodParser(rawVal string) (any, error) {
	src := strings.TrimSpace(rawVal)
	switch src {
	case "/img/vlr/game/round/elim.webp":
		return models.Eliminate, nil
	case "/img/vlr/game/round/boom.webp":
		return models.SpikeExplode, nil
	case "/img/vlr/game/round/defuse.webp":
		return models.Defuse, nil
	case "/img/vlr/game/round/time.webp":
		return models.OutOfTime, nil
	default:
		return nil, fmt.Errorf("Unable to specify the won method from this img src: %s", src)
	}
}

func teamDefParser(t1Id, t2Id int, tWonId *int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		teamWonClasses := strings.TrimSpace(rawVal)
		var teamWon, teamLost int
		if t1Id == *tWonId {
			teamWon = t1Id
			teamLost = t2Id
		} else {
			teamWon = t2Id
			teamLost = t1Id
		}

		if strings.Contains(teamWonClasses, "mod-ct") {
			return teamWon, nil
		} else if strings.Contains(teamWonClasses, "mod-t") {
			return teamLost, nil
		} else {
			return nil, fmt.Errorf("Error determine the team def from this class: %s", teamWonClasses)
		}
	}
}

// NOTE: This parse must be run after team won has been retrieved (parser for overview need to run first)
func buyTypeParser(roundNo int) htmlx.Parser {
	return func(rawVal string) (any, error) {
		buyStr := strings.TrimSpace(rawVal)

		if roundNo == 1 || roundNo == 12 {
			return models.Pistol, nil
		}

		switch buyStr {
		case "":
			return models.Eco, nil
		case "$":
			return models.SemiEco, nil
		case "$$":
			return models.SemiBuy, nil
		case "$$$":
			return models.FullBuy, nil
		default:
			return nil, fmt.Errorf("Unable to determinte the buy type from this string: %s", buyStr)
		}
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
func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	roundStats, ok := ctx.Value("roundStat").(*models.RoundStatSchema)
	if !ok {
		return fmt.Errorf("Unable to find round stats")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find the transaction")
	}

	logrus.Debug("Scraping round overview info")
	var roundOverviewSchema models.RoundOverviewSchema

	if err := htmlx.ParseFromSelection(&roundOverviewSchema, selection.Eq(0), htmlx.SetParsers(
		map[string]htmlx.Parser{
			"teamWonParser":   teamWonParser(roundStats.Team1Id, roundStats.Team2Id),
			"wonMethodParser": wonMethodParser,
			"teamDefParser":   teamDefParser(roundStats.Team1Id, roundStats.Team2Id, &roundOverviewSchema.TeamWon),
		},
	)); err != nil {
		return err
	}

	roundStats.RoundOverviewSchema = roundOverviewSchema

	logrus.Debug("Scraping round economy info")
	var roundEconomySchema models.RoundEconomySchema

	if err := htmlx.ParseFromSelection(&roundEconomySchema, selection.Eq(1), htmlx.SetParsers(
		map[string]htmlx.Parser{
			"buyTypeParser": buyTypeParser(roundStats.RoundNo),
			"balanceParser": balanceParser,
		},
	)); err != nil {
		return err
	}

	roundStats.RoundEconomySchema = roundEconomySchema

	logrus.Debug("Saving round stats to db")
	if err := tx.Table("round_stats").Create(roundStats).Error; err != nil {
		return err
	}

	return nil
}
