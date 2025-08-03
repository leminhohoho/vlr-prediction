package piper

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Backend is the interface implemented by objects those make http request and return html content.
type Backend interface {
	// Do take 2 string as the method, url and a 3rd paramter as request body to make a request of given method.
	// It return [github.com/PuerkitoBio/goquery.Selection] as a return value and a 2nd value as error.
	Do(method, url string, dat io.Reader) (*goquery.Selection, error)
	// Get take a string as the url and a 2nd paramter as request body to make a get request.
	// It return [github.com/PuerkitoBio/goquery.Selection] as a return value and a 2nd value as error.
	Get(url string, dat io.Reader) (*goquery.Selection, error)
	// Post take a string as the url and a 2nd paramter as request body to make a post request.
	// It return [github.com/PuerkitoBio/goquery.Selection] as a return value and a 2nd value as error.
	Post(url string, dat io.Reader) (*goquery.Selection, error)
	// Put take a string as the url and a 2nd paramter as request body to make a put request.
	// It return [github.com/PuerkitoBio/goquery.Selection] as a return value and a 2nd value as error.
	Put(url string, dat io.Reader) (*goquery.Selection, error)
	// Delete take a string as the url and a 2nd paramter as request body to make a delete request.
	// It return [github.com/PuerkitoBio/goquery.Selection] as a return value and a 2nd value as error.
	Delete(url string, dat io.Reader) (*goquery.Selection, error)
}

// PiperBackend use [net/http.Client] under the hood.
type PiperBackend struct {
	client *http.Client
}

// NewPiperBackend return a ready to be used [PiperBackend]
func NewPiperBackend(client *http.Client) *PiperBackend {
	return &PiperBackend{
		client: client,
	}
}

// Do implement the [Backend] interface
func (b *PiperBackend) Do(method, url string, body io.Reader) (*goquery.Selection, error) {
	if method == "" {
		method = "GET"
	}

	switch method {
	case "GET", "POST", "PUT", "DELETE":
		break
	default:
		return nil, fmt.Errorf("Invalid method '%s'", method)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc.Selection, nil
}

// Get implement the [Backend] interface
func (b *PiperBackend) Get(url string, body io.Reader) (*goquery.Selection, error) {
	return b.Do("GET", url, body)
}

// Post implement the [Backend] interface
func (b *PiperBackend) Post(url string, body io.Reader) (*goquery.Selection, error) {
	return b.Do("POST", url, body)
}

// Put implement the [Backend] interface
func (b *PiperBackend) Put(url string, body io.Reader) (*goquery.Selection, error) {
	return b.Do("PUT", url, body)
}

// Delete implement the [Backend] interface
func (b *PiperBackend) Delete(url string, body io.Reader) (*goquery.Selection, error) {
	return b.Do("DELETE", url, body)
}
