package crawler

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/customerrors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const (
	initializeScript = `
				DROP TABLE IF EXISTS matches_to_be_scraped;
				CREATE TABLE IF NOT EXISTS matches_to_be_scraped (
    				id INTEGER PRIMARY KEY,
    				url TEXT UNIQUE NOT NULL,
    				date TEXT NOT NULL,
    				failed INTEGER CHECK(failed IN (0, 1)) DEFAULT 0,
    				error TEXT
				);
				`
	dateSelector  = `#wrapper > div.col-container > div > div.wf-label.mod-large`
	matchSelector = `a[href].match-item`
	dateLayout    = "Mon, January 2, 2006"
)

type MatchToBeScraped struct {
	Url  string
	Date time.Time
}

// NOTE: createCacheDb() and initializeCacheDb() only create and initialize the temp db if it is not exist/not initialized

func createCacheDb(tmpCacheDbPath string) (*sql.DB, error) {
	_, err := os.Stat(tmpCacheDbPath)
	if err == nil {
		logrus.Debug("cache db has been create, continue the process")
		db, err := sql.Open("sqlite3", tmpCacheDbPath)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	logrus.Debug("cache db is not created, create a new one")
	if _, err = os.Create(tmpCacheDbPath); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", tmpCacheDbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initializeScript)

	return db, nil
}

func isDbEmpty(db *sql.DB, tableName string) (bool, error) {
	var count int

	row := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
	if err := row.Scan(&count); err != nil {
		return true, err
	}

	return count == 0, nil
}

func getLatestDate(db *sql.DB, tableName string) (time.Time, error) {
	var dateLimitStr string

	row := db.QueryRow(fmt.Sprintf("SELECT date FROM %s ORDER BY date DESC LIMIT 1", tableName))
	if err := row.Scan(&dateLimitStr); err != nil {
		return time.Time{}, err
	}

	date, err := time.Parse("2006-01-02 15:04:05-07:00", dateLimitStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func crawlMatchesPerDate(
	matchesContainer *goquery.Selection,
	matchDate time.Time,
) ([]MatchToBeScraped, error) {
	errChan := make(chan error)
	resChan := make(chan []MatchToBeScraped)

	go func() {
		var matchesToBeScraped []MatchToBeScraped

		matchNodes := matchesContainer.Find(matchSelector)
		if matchNodes.Length() == 0 {
			errChan <- customerrors.ErrMissingHTMLSelection{Doc: matchesContainer}
			return
		}

		matchesContainer.Find(matchSelector).Each(func(_ int, matchNode *goquery.Selection) {
			matchUrl, exists := matchNode.Attr("href")
			if !exists {
				html, err := matchNode.Html()
				if err != nil {
					errChan <- err
					return
				}

				errChan <- fmt.Errorf("Error getting link from node:\n%s", html)
				return
			}

			logrus.WithFields(logrus.Fields{
				"url":  matchUrl,
				"date": matchDate.Format(dateLayout),
			}).Info("Match to scrape info")

			matchesToBeScraped = append(matchesToBeScraped, MatchToBeScraped{
				Url:  matchUrl,
				Date: matchDate,
			})
		})

		resChan <- matchesToBeScraped
	}()

	select {
	case err := <-errChan:
		return nil, err
	case res := <-resChan:
		return res, nil
	}
}

func crawlMatchesUpToDate(dateLimit time.Time) ([]MatchToBeScraped, error) {
	errChan := make(chan error)
	resChan := make(chan []MatchToBeScraped)

	go func() {
		pageCount := 1
		var matchesToBeScraped []MatchToBeScraped

		crawler := colly.NewCollector(colly.AllowedDomains("www.vlr.gg"))

		crawler.OnRequest(func(req *colly.Request) {
			logrus.Debugf("Crawler visiting: %s\n", req.AbsoluteURL(req.URL.String()))
		})

		crawler.OnResponse(func(res *colly.Response) {
			logrus.Debugf("Finding matches on page %d", pageCount)
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body))
			if err != nil {
				errChan <- err
				return
			}

			doc.Find(dateSelector).Each(func(_ int, dateNode *goquery.Selection) {
				dateStr := strings.TrimSpace(dateNode.Children().Remove().End().Text())
				matchDate, err := time.Parse(dateLayout, dateStr)
				if err != nil {
					errChan <- err
					return
				}

				if matchDate.Before(dateLimit) || matchDate.Equal(dateLimit) {
					resChan <- matchesToBeScraped
					return
				}

				matchesContainer := dateNode.Next()
				if matchesContainer.Length() == 0 {
					errChan <- customerrors.ErrMissingHTMLSelection{Doc: dateNode.Parent()}
				}

				logrus.Debugf("Finding matches on %s", dateStr)
				currentPageMatches, err := crawlMatchesPerDate(matchesContainer, matchDate)
				if err != nil {
					errChan <- err
					return
				}

				matchesToBeScraped = append(matchesToBeScraped, currentPageMatches...)
			})

			pageCount++
			crawler.Visit(fmt.Sprintf("https://www.vlr.gg/matches/results/?page=%d", pageCount))
		})

		crawler.OnError(func(res *colly.Response, err error) {
			errChan <- err
		})

		crawler.Visit(fmt.Sprintf("https://www.vlr.gg/matches/results/?page=%d", pageCount))
	}()

	select {
	case err := <-errChan:
		return nil, err
	case res := <-resChan:
		logrus.Infof("Number of new matches to be scraped: %d", len(res))
		return res, nil
	}
}

func getDateLimit(tmpCacheDb, vlrDb *sql.DB) (time.Time, error) {
	tmpCacheDbEmpty, err := isDbEmpty(tmpCacheDb, "matches_to_be_scraped")
	if err != nil {
		return time.Time{}, err
	}

	vlrDbEmpty, err := isDbEmpty(vlrDb, "matches")
	if err != nil {
		return time.Time{}, err
	}

	var dateLimit time.Time

	if !tmpCacheDbEmpty {
		dateLimit, err = getLatestDate(tmpCacheDb, "matches_to_be_scraped")
		if err != nil {
			return time.Time{}, err
		}
	} else if !vlrDbEmpty {
		dateLimit, err = getLatestDate(vlrDb, "matches")
		if err != nil {
			return time.Time{}, err
		}
	} else {
		dateLimit, err = time.Parse("2006-01-02", os.Getenv("DATE_LIMIT"))
		if err != nil {
			return time.Time{}, err
		}
	}

	return dateLimit, nil
}

func CrawlMatches(tmpCacheDbPath, vlrDbPath string) ([]MatchToBeScraped, error) {
	var err error

	logrus.Debug("Connecting to cache db")
	tmpCacheDb, err := createCacheDb(tmpCacheDbPath)
	if err != nil {
		return nil, err
	}
	defer tmpCacheDb.Close()

	logrus.Debug("Connecting to vlr db")
	vlrDb, err := sql.Open("sqlite3", vlrDbPath)
	if err != nil {
		return nil, err
	}
	defer vlrDb.Close()

	logrus.Debug("Determine the limit date where matches haven't been scraped")
	dateLimit, err := getDateLimit(tmpCacheDb, vlrDb)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Start scraping matches after %s", dateLimit.Format(dateLayout))
	matchesToBeScraped, err := crawlMatchesUpToDate(dateLimit)
	if err != nil {
		return nil, err
	}

	logrus.Debug("Insert the matches to cache db")
	for _, matchToBeScraped := range matchesToBeScraped {
		if _, err = tmpCacheDb.Exec(
			"INSERT INTO matches_to_be_scraped(url, date) VALUES(?,?)",
			matchToBeScraped.Url, matchToBeScraped.Date,
		); err != nil {
			return nil, err
		}
	}

	return matchesToBeScraped, nil
}
