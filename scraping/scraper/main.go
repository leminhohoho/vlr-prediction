package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/piper"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/crawler"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/matches"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/matchmaps"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/playerduelstats"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/playerstats"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/roundstats"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/teams"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/scrapers/tournaments"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/urlinfo"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.TraceLevel)

	vlrDb, err := gorm.Open(sqlite.Open(os.Getenv("VLR_DB_PATH")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	cacheDb, err := sql.Open("sqlite3", path.Join(os.Getenv("TMP_DIR"), "vlr_cache.db"))
	if err != nil {
		panic(err)
	}

	cache, err := piper.NewCacheDb(path.Join(os.Getenv("TMP_DIR"), "scraper_cache.db"))
	if err != nil {
		panic(err)
	}

	if err = cache.Validate(); err != nil && err != piper.ErrIncorrectSchema {
		panic(err)
	} else if err == piper.ErrIncorrectSchema {
		if err = cache.Setup(); err != nil {
			panic(err)
		}
	}

	backend := piper.NewPiperBackend(&http.Client{})

	sc := piper.NewScraper(backend, cache)
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/[0-9]+\/[a-z0-9\/-]*$`), matches.Handler)
	sc.Handle(regexp.MustCompile(`matchMaps`), matchmaps.Handler)
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/team\/[0-9]+\/[a-z0-9\/-]*$`), teams.Handler)
	sc.Handle(regexp.MustCompile(`^https:\/\/www\.vlr\.gg\/event\/[0-9]+\/[a-z0-9\/-]*$`), tournaments.Handler)
	sc.Handle(regexp.MustCompile(`^roundStat$`), roundstats.Handler)
	sc.Handle(regexp.MustCompile(`^playerStats$`), playerstats.Handler)
	sc.Handle(regexp.MustCompile(`^duelStats$`), playerduelstats.Handler)

	matchesToBeScraped, err := crawler.CrawlMatches(
		path.Join(os.Getenv("TMP_DIR"), "vlr_cache.db"),
		os.Getenv("VLR_DB_PATH"),
	)
	if err != nil {
		panic(err)
	}

	var errs []error

	for i, matchToBeScraped := range matchesToBeScraped {
		urlInfo, err := urlinfo.ExtractUrlInfo(matchToBeScraped.Url)
		if err != nil {
			logrus.Errorf("Unable to extraction information from url, skip to next match")
			continue
		}

		fullUrl := "https://www.vlr.gg" + matchToBeScraped.Url

		var exists bool

		if err := vlrDb.Table("matches").Select("count(*) > 0").Where("id = ?", urlInfo.Id).Find(&exists).Error; err != nil {
			logrus.Fatal(err)
		} else if exists {
			logrus.Debug("Match exists, continue")
			continue
		}

		if i%50 == 0 && i > 0 {
			logrus.Debug("Pause for 30 second")
			time.Sleep(time.Second * 30)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logrus.Debugf("Scraping from: %s", "https://www.vlr.gg"+matchToBeScraped.Url)

		matchSchema := models.MatchSchema{Id: urlInfo.Id, Url: fullUrl, Date: matchToBeScraped.Date}

		if err := vlrDb.Transaction(func(tx *gorm.DB) error {
			ctx = context.WithValue(context.WithValue(ctx, "matchSchema", &matchSchema), "tx", tx)

			overviewRes, err := http.Get(fullUrl + "/?games=all&tab=overview")
			if err != nil {
				return fmt.Errorf("Error fetching from overview page, skip to next match")
			}

			performanceRes, err := http.Get(fullUrl + "/?games=all&tab=performance")
			if err != nil {
				return fmt.Errorf("Error fetching from performance page, skip to next match")
			}

			economyRes, err := http.Get(fullUrl + "/?games=all&tab=economy")
			if err != nil {
				return fmt.Errorf("Error fetching from economy page, skip to next match")
			}

			overviewDoc, err := goquery.NewDocumentFromReader(overviewRes.Body)
			if err != nil {
				return fmt.Errorf("Error parsing overview content: %s", err.Error())
			}

			performanceDoc, err := goquery.NewDocumentFromReader(performanceRes.Body)
			if err != nil {
				return fmt.Errorf("Error parsing performance content: %s", err.Error())
			}

			economyDoc, err := goquery.NewDocumentFromReader(economyRes.Body)
			if err != nil {
				return fmt.Errorf("Error parsing economy content: %s", err.Error())
			}

			combined := overviewDoc.Selection.AddSelection(performanceDoc.Selection).AddSelection(economyDoc.Selection)

			if err := sc.Pipe(fullUrl, ctx, combined); err != nil {
				return fmt.Errorf("Error: '%s', skip to next match", err.Error())
			}

			return nil
		}); err != nil {
			logrus.Error(err)
			errs = append(errs, fmt.Errorf("Error scraping from '%s': %s", fullUrl, err.Error()))
			continue
		}

		if _, err := cacheDb.Exec("DELETE FROM matches_to_be_scraped WHERE url = ?", matchToBeScraped.Url); err != nil {
			logrus.Errorf("Error deleting match from cache: '%s', skip to next match", err.Error())
		}
	}

	fmt.Println("=========================== ERROR ===========================")
	for _, err := range errs {
		logrus.Error(err)
	}
}
