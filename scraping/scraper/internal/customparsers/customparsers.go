package customparsers

import (
	"strings"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/urlinfo"
)

func IdParser(rawVal string) (any, error) {
	url := strings.TrimSpace(rawVal)
	vlrUrlInfo, err := urlinfo.ExtractUrlInfo(url)
	if err != nil {
		return nil, err
	}

	return vlrUrlInfo.Id, nil
}
