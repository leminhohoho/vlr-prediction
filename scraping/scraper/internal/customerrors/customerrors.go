package customerrors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// Error for missing HTML elements
type ErrMissingHTMLSelection struct {
	Doc *goquery.Selection
}

func (e ErrMissingHTMLSelection) Error() string {
	html, err := e.Doc.Html()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Error: empty goquery selection, parent Node:\n%s", html)
}
