package matches

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	matchMapSelector        = `#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div[data-game-id="%s"]:has(div+div)`
	matchMapGenericSelector = `#wrapper > div.col-container > div.col.mod-3 > div:nth-child(6) > div > div.vm-stats-container > div[data-game-id!="all"]:has(div+div)`
)

func stageParser(rawVal string) (any, error) {
	matchHeader := helpers.ToSnakeCase(strings.TrimSpace(rawVal))
	if strings.Contains(matchHeader, "grand_final") {
		return models.GrandFinal, nil
	}

	if strings.Contains(matchHeader, "playoff") {
		return models.Playoff, nil
	}

	return models.GroupStage, nil
}

func ratingParser(rawVal string) (any, error) {
	ratingStr := strings.TrimSpace(rawVal)
	if ratingStr == "" {
		return 0, nil
	}

	if !regexp.MustCompile(`^\[[0-9]+\]$`).MatchString(ratingStr) {
		return -1, fmt.Errorf("rating string %s is not valid for converting to rating", ratingStr)
	}

	ratingStr = ratingStr[1 : len(ratingStr)-1]
	rating, _ := strconv.Atoi(ratingStr)
	return rating, nil
}

func Handler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	matchSchema, ok := ctx.Value("matchSchema").(*models.MatchSchema)
	if !ok {
		return fmt.Errorf("Unable to find match schema")
	}

	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find gorm transaction")
	}

	overviewContent := selection.Eq(0)
	performanceContent := selection.Eq(1)
	economyContent := selection.Eq(2)

	parsers := map[string]htmlx.Parser{
		"idParser":     customparsers.IdParser,
		"stageParser":  stageParser,
		"ratingParser": ratingParser,
	}

	logrus.Debug("Parsing information from html onto match schema")
	if err := htmlx.ParseFromSelection(matchSchema, overviewContent, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	jsonDat, err := json.MarshalIndent(*matchSchema, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonDat))

	logrus.Debug("Saving match to db")
	if err := tx.Table("matches").Create(matchSchema).Error; err != nil {
		return err
	}

	logrus.Debug("Locating maps nodes")

	errChan := make(chan error)

	var mu sync.Mutex

	go func() {
		overviewContent.Find(matchMapGenericSelector).Each(func(_ int, mapOverviewNode *goquery.Selection) {
			mu.Lock()

			gameId, exists := mapOverviewNode.Attr("data-game-id")
			if !exists {
				errChan <- fmt.Errorf("Unable to find game id")
				html, _ := mapOverviewNode.Html()
				fmt.Println(html)
				return
			}

			mapPerformanceNode := performanceContent.Find(fmt.Sprintf(matchMapSelector, gameId))
			mapEconomyNode := economyContent.Find(fmt.Sprintf(matchMapSelector, gameId))

			combined := mapOverviewNode.AddSelection(mapPerformanceNode).AddSelection(mapEconomyNode)

			matchMap := models.MatchMapSchema{
				MatchId: matchSchema.Id,
				Team1Id: matchSchema.Team1Id,
				Team2Id: matchSchema.Team2Id,
			}

			ctx := context.WithValue(context.WithValue(context.Background(), "matchMapSchema", &matchMap), "tx", tx)

			if err := sc.Pipe("matchMaps", ctx, combined); err != nil {
				errChan <- err
				return
			}

			mu.Unlock()
		})

		errChan <- nil
	}()

	select {
	case err := <-errChan:
		return err
	}
}
