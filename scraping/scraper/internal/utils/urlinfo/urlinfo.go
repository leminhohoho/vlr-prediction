package urlinfo

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type VlrDomain int

const (
	Matches VlrDomain = iota
	Teams
	Tournaments
	Players

	matchUrlPattern  = `^\/[0-9]+\/[a-zA-Z0-9-]+$`
	playerUrlPattern = `^\/player\/[0-9]+\/[a-zA-Z0-9-]+$`
	teamUrlPattern   = `^\/team\/[0-9]+\/[a-zA-Z0-9-]+$`
	eventUrlPattern  = `^\/event\/[0-9]+\/[a-zA-Z0-9-]+\/?[a-zA-Z0-9-]*$`
)

func getMatchId(urlStr string) int {
	idStr := strings.Split(urlStr, "/")[1]
	id, _ := strconv.Atoi(idStr)
	return id
}

func getTeamId(urlStr string) int {
	idStr := strings.Split(urlStr, "/")[2]
	id, _ := strconv.Atoi(idStr)
	return id
}

func getPlayerId(urlStr string) int {
	idStr := strings.Split(urlStr, "/")[2]
	id, _ := strconv.Atoi(idStr)
	return id
}

func getTournamentId(urlStr string) int {
	idStr := strings.Split(urlStr, "/")[2]
	id, _ := strconv.Atoi(idStr)
	return id
}

type VlrUrl struct {
	UrlType VlrDomain
	Id      int
	Url     *url.URL
}

func (v VlrUrl) IsMatch() bool {
	return v.UrlType == Matches
}

func (v VlrUrl) IsTeam() bool {
	return v.UrlType == Teams
}

func (v VlrUrl) IsPlayer() bool {
	return v.UrlType == Players
}

func (v VlrUrl) IsTournament() bool {
	return v.UrlType == Tournaments
}

func ExtractUrlInfo(urlStr string) (VlrUrl, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return VlrUrl{}, err
	}

	relativeUrl := parsedUrl.Path

	if regexp.MustCompile(matchUrlPattern).MatchString(relativeUrl) {
		return VlrUrl{
			UrlType: Matches,
			Id:      getMatchId(relativeUrl),
			Url:     parsedUrl,
		}, nil
	} else if regexp.MustCompile(playerUrlPattern).MatchString(relativeUrl) {
		return VlrUrl{
			UrlType: Players,
			Id:      getPlayerId(relativeUrl),
			Url:     parsedUrl,
		}, nil
	} else if regexp.MustCompile(teamUrlPattern).MatchString(relativeUrl) {
		return VlrUrl{
			UrlType: Teams,
			Id:      getTeamId(relativeUrl),
			Url:     parsedUrl,
		}, nil
	} else if regexp.MustCompile(eventUrlPattern).MatchString(relativeUrl) {
		return VlrUrl{
			UrlType: Tournaments,
			Id:      getTournamentId(relativeUrl),
			Url:     parsedUrl,
		}, nil
	} else {
		return VlrUrl{}, fmt.Errorf("Url is not from www.vlr.gg domain")
	}
}
