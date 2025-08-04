package matches

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customparsers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/helpers"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func MatchHandler(sc *piper.Scraper, ctx context.Context, selection *goquery.Selection) error {
	logrus.Debug("Parsing information from match html content into match schema")
	matchSchema, ok := ctx.Value("matchSchema").(*models.MatchSchema)
	if !ok {
		return fmt.Errorf("Unable to find match schema")
	}

	_, ok = ctx.Value("tx").(*gorm.DB)
	if !ok {
		return fmt.Errorf("Unable to find gorm transaction")
	}

	parsers := map[string]htmlx.Parser{
		"idParser":     customparsers.IdParser,
		"stageParser":  stageParser,
		"ratingParser": ratingParser,
	}

	if err := htmlx.ParseFromSelection(matchSchema, selection, htmlx.SetParsers(parsers)); err != nil {
		return err
	}

	jsonDat, err := json.MarshalIndent(*matchSchema, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonDat))

	return nil
}
