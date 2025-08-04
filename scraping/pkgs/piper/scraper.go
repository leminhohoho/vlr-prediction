package piper

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Handler represent a function that is invoke the scrape from a [github.com/PuerkitoBio/goquery.Selection].
// It take 1st agruement as [Scraper] and 2nd one as the selection to scrape from.
type Handler func(sc *Scraper, ctx context.Context, selection *goquery.Selection) error

// Scraper represent a distributed web scraping server.
// All scraper methods are safe for concurrent usage.
type Scraper struct {
	mu sync.Mutex

	handlers map[*regexp.Regexp]Handler
	backend  Backend
	cache    Cache
}

// NewScraper create a new scraper with the given backend and cache.
func NewScraper(backend Backend, cache Cache) *Scraper {
	return &Scraper{
		backend:  backend,
		cache:    cache,
		handlers: map[*regexp.Regexp]Handler{},
	}
}

// Handle take a regex expression with a [Handler] and register a handler for the given pattern.
// Handle is safe for concurrent usage.
func (sc *Scraper) Handle(regex *regexp.Regexp, h Handler) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.handlers[regex] = h
}

// Pipe calls the handler that associated with the first regex expression that match the pattern.
// It return error if either the no matching regex expression is found or the handler return error.
// Pipe is safe for concurrent usage.
func (sc *Scraper) Pipe(pattern string, ctx context.Context, selection *goquery.Selection) error {
	for regex, handler := range sc.handlers {
		if regex.MatchString(pattern) {
			return handler(sc, ctx, selection)
		}
	}

	return fmt.Errorf("No handler match the pattern: '%s'", pattern)
}

// Get make a GET request using the given url and body and pass the [github.com/PuerkitoBio/goquery.Selection] to [Scraper.Pipe]
func (sc *Scraper) Get(url string, ctx context.Context, body io.Reader) error {
	selection, err := sc.backend.Get(url, body)
	if err != nil {
		return err
	}

	return sc.Pipe(url, ctx, selection)
}

// Post make a POST request using the given url and body and pass the [github.com/PuerkitoBio/goquery.Selection] to [Scraper.Pipe]
func (sc *Scraper) Post(url string, ctx context.Context, body io.Reader) error {
	selection, err := sc.backend.Post(url, body)
	if err != nil {
		return err
	}

	return sc.Pipe(url, ctx, selection)
}

// Put make a PUT request using the given url and body and pass the [github.com/PuerkitoBio/goquery.Selection] to [Scraper.Pipe]
func (sc *Scraper) Put(url string, ctx context.Context, body io.Reader) error {
	selection, err := sc.backend.Put(url, body)
	if err != nil {
		return err
	}

	return sc.Pipe(url, ctx, selection)
}

// Delete make a DELETE request using the given url and body and pass the [github.com/PuerkitoBio/goquery.Selection] to [Scraper.Pipe]
func (sc *Scraper) Delete(url string, ctx context.Context, body io.Reader) error {
	selection, err := sc.backend.Delete(url, body)
	if err != nil {
		return err
	}

	return sc.Pipe(url, ctx, selection)
}

// Cache return the cache implementation used by the scraper
func (sc *Scraper) Cache() Cache {
	return sc.cache
}
