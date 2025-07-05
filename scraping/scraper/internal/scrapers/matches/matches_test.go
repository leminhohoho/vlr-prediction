package matches

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func setUpLogrus(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.Warnf("Invalid log level '%s', defaulting to info", os.Getenv("LOG_LEVEL"))
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
}

func TestMatchWithFullInformation(t *testing.T) {
	setUpLogrus(t)

	res, err := http.Get("https://www.vlr.gg/506931/bilibili-gaming-vs-tyloo-vct-2025-china-stage-2-w1")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	matchScraper := NewMatchScraper(
		nil,
		nil,
		doc.Selection,
		506931,
		"https://www.vlr.gg/506931/bilibili-gaming-vs-tyloo-vct-2025-china-stage-2-w1",
		time.Now(),
	)

	if err := matchScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	matchScraper.PrettyPrint()
}

func TestMatchWithMissingInformation(t *testing.T) {
	setUpLogrus(t)

	res, err := http.Get(
		"https://www.vlr.gg/489315/twisted-minds-vs-villianarc-challengers-2025-mena-resilience-gcc-pakistan-iraq-split-2-w5",
	)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	matchScraper := NewMatchScraper(
		nil,
		nil,
		doc.Selection,
		489315,
		"https://www.vlr.gg/489315/twisted-minds-vs-villianarc-challengers-2025-mena-resilience-gcc-pakistan-iraq-split-2-w5",
		time.Now(),
	)

	if err := matchScraper.Scrape(); err != nil {
		t.Fatal(err)
	}

	matchScraper.PrettyPrint()
}
