package customparsers

import (
	"strings"

	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/utils/urlinfo"
	"gorm.io/gorm"
)

func IdParser(rawVal string) (any, error) {
	url := strings.TrimSpace(rawVal)
	vlrUrlInfo, err := urlinfo.ExtractUrlInfo(url)
	if err != nil {
		return nil, err
	}

	return vlrUrlInfo.Id, nil
}

func MapIdParser(tx *gorm.DB) htmlx.Parser {
	return func(rawVal string) (any, error) {
		mapName := strings.TrimSpace(rawVal)
		var vlrMap models.MapSchema

		rs := tx.Table("maps").Where("name = ?", mapName).First(&vlrMap)
		if rs.Error != nil {
			return -1, rs.Error
		}

		return vlrMap.Id, nil
	}
}
