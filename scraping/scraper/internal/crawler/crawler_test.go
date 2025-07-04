package crawler

import (
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func TestCrawler(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logrus.Warnf("Invalid log level '%s', defaulting to info", os.Getenv("LOG_LEVEL"))
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)

	_, err = CrawlMatches(
		path.Join(os.Getenv("TMP_DIR"), "vlr_cache.db"),
		"../../../../database/vlr.db",
	)
	if err != nil {
		t.Fatal(err)
	}
}
