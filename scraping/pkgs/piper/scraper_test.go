package piper

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestScraper(t *testing.T) {
	sqliteDbPath := "tmp/test_cache.db"

	cache, err := NewCacheDb(sqliteDbPath)
	if err != nil {
		t.Fatal(err)
	}

	if err = cache.Validate(); err != nil && err != ErrIncorrectSchema {
		t.Fatal(err)
	} else if err == ErrIncorrectSchema {
		if err = cache.Setup(); err != nil {
			t.Fatal(err)
		}
	}

	backend := NewPiperBackend(&http.Client{})

	s := NewScraper(backend, cache)

	s.Handle(
		regexp.MustCompile(`^https:\/\/books\.toscrape\.com\/$`),
		func(sc *Scraper, ctx context.Context, selection *goquery.Selection) error {
			categories := selection.Find(`#default > div > div > div > aside > div.side_categories > ul > li > ul`)
			fmt.Printf("Number of categories: %d\n", categories.Children().Length())
			categories.Children().Each(func(_ int, category *goquery.Selection) {
				fmt.Printf("	- %s\n", strings.TrimSpace(category.Find("a").Text()))

				url, exists := category.Find("a").Attr("href")
				if !exists {
					t.Errorf("Url does not exist")
				}

				sc.Get("https://books.toscrape.com/"+url, ctx, nil)
			})

			return nil
		},
	)

	s.Handle(
		regexp.MustCompile(`^https://books\.toscrape\.com/catalogue/category/[a-z0-9\/_.]+$`),
		func(sc *Scraper, ctx context.Context, selection *goquery.Selection) error {
			numOfResults, err := strconv.Atoi(
				strings.TrimSpace(selection.Find(`#default > div > div > div > div > form > strong`).Text()),
			)
			if err != nil {
				t.Error(err)
			}

			fmt.Printf("		+ Number of books: %d\n", numOfResults)

			return nil
		},
	)

	s.Get("https://books.toscrape.com/", context.Background(), nil)
}
