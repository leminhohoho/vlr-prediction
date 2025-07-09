package players

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlayerSchema struct {
	Id        int
	Name      string  `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h1"`
	RealName  *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h2" parser:"realNameParser"`
	Url       string
	ImgUrl    *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div.wf-avatar.mod-player > div > img"     parser:"imgUrlParser"    source:"attr=src"`
	CountryId *int    `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div.ge-text-light"     parser:"countryIdParser"`
}

type PlayerScraper struct {
	Data              PlayerSchema
	PlayerPageContent *goquery.Selection
	Conn              *sql.DB
	Tx                *gorm.Tx
}

func NewPlayerScraper(
	conn *sql.DB,
	tx *gorm.Tx,
	playerPageContent *goquery.Selection,
	playerId int,
	playerUrl string,
) *PlayerScraper {
	var realName, imgUrl string
	var countryId int

	return &PlayerScraper{
		Data: PlayerSchema{
			Id:        playerId,
			Url:       playerUrl,
			RealName:  &realName,
			ImgUrl:    &imgUrl,
			CountryId: &countryId,
		},
		PlayerPageContent: playerPageContent,
		Tx:                tx,
		Conn:              conn,
	}
}

func realNameParser(rawVal string) (any, error) {
	nameStr := strings.TrimSpace(rawVal)

	if nameStr == "" {
		return nil, nil
	}

	return &nameStr, nil
}

func imgUrlParser(rawVal string) (any, error) {
	imgUrlStr := strings.TrimSpace(rawVal)

	if imgUrlStr == "" {
		return nil, nil
	}

	return &imgUrlStr, nil
}

func getCountryId(conn *sql.DB, countryName string) (int, error) {
	var id int

	row := conn.QueryRow("SELECT id FROM countries WHERE name = ?", countryName)

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (p *PlayerScraper) countryIdParser(rawVal string) (any, error) {
	countryName := strings.TrimSpace(rawVal)

	if countryName == "" {
		return nil, nil
	}

	fmt.Printf("Player country: %s\n", countryName)

	countryId, err := getCountryId(p.Conn, countryName)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}

		// TODO: Add code for inserting into (the line below is placeholder)
		return nil, nil
	}

	return &countryId, nil
}

func (p *PlayerScraper) PrettyPrint() error {
	jsonStr, err := json.MarshalIndent(p.Data, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonStr))

	return nil
}

func (p *PlayerScraper) Scrape() error {
	parsers := map[string]htmlx.Parser{
		"realNameParser":  realNameParser,
		"imgUrlParser":    imgUrlParser,
		"countryIdParser": p.countryIdParser,
	}

	logrus.Debug("Scraping player information")
	if err := htmlx.ParseFromSelection(&p.Data, p.PlayerPageContent, htmlx.SetParsers(parsers), htmlx.SetAllowNilPointer(true)); err != nil {
		return err
	}

	return nil
}
