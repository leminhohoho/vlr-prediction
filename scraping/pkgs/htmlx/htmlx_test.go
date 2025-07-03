package htmlx

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://www.vlr.gg/498628/paper-rex-vs-fnatic-valorant-masters-toronto-2025-gf/?game=221168&tab=economy"

type BanPickLog struct {
	Value string `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-note"`
}

type MatchInfo struct {
	Team1Name     string `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1 > div > div.wf-title-med"`
	Team2Name     string `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2 > div > div.wf-title-med.mod-single"`
	Team1Score    int    `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-winner"`
	Team2Score    int    `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span.match-header-vs-score-loser"`
	TournamentURL string `selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.mod-color.mod-bg-after-striped_purple.match-header > div.match-header-super > div:nth-child(1) > a"                                                                source:"attr=href"`
	TeamWonBet    int    `selector:"#wrapper > div.col-container > div.col.mod-3 > div:nth-child(2) > a:nth-child(2) > div > div.match-bet-item-team > span:nth-child(4)"`
	BanPickLog    BanPickLog
}

func TestHTMLx(t *testing.T) {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	matchInfo := MatchInfo{}

	if err = ParseFromDocument(&matchInfo, doc); err != nil {
		t.Fatal(err)
	}

	if matchInfo.Team1Name != "Paper Rex" {
		t.Errorf("Wrong team 1 name, want 'Paper Rex', get '%s'", matchInfo.Team1Name)
	}

	if matchInfo.Team2Name != "FNATIC" {
		t.Errorf("Wrong team 2 name, want 'FNATIC', get '%s'", matchInfo.Team2Name)
	}

	if matchInfo.Team1Score != 3 {
		t.Errorf("Wrong team 1 score, want 3, get %d", matchInfo.Team1Score)
	}

	if matchInfo.Team2Score != 1 {
		t.Errorf("Wrong team 2 score, want 1, get %d", matchInfo.Team1Score)
	}

	if matchInfo.TournamentURL != "/event/2282/valorant-masters-toronto-2025/playoffs" {
		t.Errorf(
			"Wrong tournament url, want '/event/2282/valorant-masters-toronto-2025/playoffs', get '%s'",
			matchInfo.TournamentURL,
		)
	}

	if matchInfo.TeamWonBet != 179 {
		t.Errorf("Wrong team won bet ammount, want 179, get %d", matchInfo.TeamWonBet)
	}

	if matchInfo.BanPickLog.Value != "PRX ban Haven; PRX ban Ascent; PRX pick Sunset; FNC pick Icebox; PRX pick Pearl; FNC pick Lotus; Split remains" {
		t.Errorf(
			"Wrong team ban pick log, want 'PRX ban Haven; PRX ban Ascent; PRX pick Sunset; FNC pick Icebox; PRX pick Pearl; FNC pick Lotus; Split remains', get '%s'",
			matchInfo.BanPickLog.Value,
		)
	}

	fmt.Println(matchInfo.Team1Name)
	fmt.Println(matchInfo.Team2Name)
	fmt.Println(matchInfo.Team1Score)
	fmt.Println(matchInfo.Team2Score)
	fmt.Println(matchInfo.TournamentURL)
	fmt.Println(matchInfo.TeamWonBet)
	fmt.Println(matchInfo.BanPickLog.Value)
}
