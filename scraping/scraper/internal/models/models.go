package models

type Side string
type AgentType string

const (
	Def Side = "def"
	Atk Side = "atk"

	Duelist    AgentType = "duelist"
	Controller AgentType = "controller"
	Sentinel   AgentType = "sentinal"
	Initiator  AgentType = "initiator"
)

type CountrySchema struct {
	Id       int    `gorm:"column:id;primaryKey;autoIcrement"`
	Name     string `gorm:"column:name"`
	RegionId int    `gorm:"column:region_id"`
}

type RegionSchema struct {
	Id   int    `gorm:"column:id;primaryKey;autoIcrement"`
	Name string `gorm:"column:name"`
}

type AgentSchema struct {
	Id          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name"`
	AgentType   AgentType `gorm:"column:agent_type"`
	ReleaseDate string    `gorm:"column:release_date"`
}

type PlayerSchema struct {
	Id        int
	Name      string  `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h1"`
	RealName  *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h2"`
	Url       string
	ImgUrl    *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div.wf-avatar.mod-player > div > img"     source:"attr=src"`
	CountryId *int    `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div.ge-text-light"                       parser:"countryIdParser"`
}

type PlayerOverviewStatSchema struct {
	MatchId     int
	MapId       int
	TeamId      int
	Side        Side
	PlayerId    int      `selector:"td.mod-player > div > a"                           source:"attr=href"  parser:"playerIdParser"`
	AgentId     int      `selector:"td.mod-agents > div > span > img"                  source:"attr=title" parser:"agentParser"`
	Rating      *float64 `selector:"td:nth-child(3) > span > span"`
	Acs         *float64 `selector:"td:nth-child(4) > span > span"`
	Kills       *int     `selector:"td:nth-child(5) > span > span"`
	Deaths      *int     `selector:"td:nth-child(6) > span > span:nth-child(2) > span"`
	Assists     *int     `selector:"td:nth-child(7) > span > span"`
	Kast        *float64 `selector:"td:nth-child(9) > span > span"`
	Adr         *float64 `selector:"td:nth-child(10) > span > span"`
	Hs          *float64 `selector:"td:nth-child(11) > span > span"`
	FirstKills  *int     `selector:"td:nth-child(12) > span > span"`
	FirstDeaths *int     `selector:"td:nth-child(13) > span > span"`
}
