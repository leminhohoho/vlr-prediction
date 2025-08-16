package models

import "time"

type Stage string
type Side string
type AgentType string
type VetoAction string
type HighlightType string
type BuyType string
type WonMethod string

const (
	Def Side = "def"
	Atk Side = "atk"

	Duelist    AgentType = "duelist"
	Controller AgentType = "controller"
	Sentinel   AgentType = "sentinal"
	Initiator  AgentType = "initiator"

	GroupStage Stage = "group_stage"
	Playoff    Stage = "playoff"
	GrandFinal Stage = "grand_final"

	BanMap    VetoAction = "ban"
	PickMap   VetoAction = "pick"
	RemainMap VetoAction = "remains"

	P2k  HighlightType = "2k"
	P3k  HighlightType = "3k"
	P4k  HighlightType = "4k"
	P5k  HighlightType = "5k"
	P1v1 HighlightType = "1v1"
	P1v2 HighlightType = "1v2"
	P1v3 HighlightType = "1v3"
	P1v4 HighlightType = "1v4"
	P1v5 HighlightType = "1v5"

	Pistol  BuyType = "pistol"
	Eco     BuyType = "eco"
	SemiEco BuyType = "semi_eco"
	SemiBuy BuyType = "semi_buy"
	FullBuy BuyType = "full_buy"

	Eliminate    WonMethod = "eliminate"
	SpikeExplode WonMethod = "spike_explode"
	Defuse       WonMethod = "defuse"
	OutOfTime    WonMethod = "out_of_time"
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

type MapSchema struct {
	Id          int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `gorm:"column:name"`
	ReleaseDate string `gorm:"column:release_date"`
}

type MatchSchema struct {
	Id           int
	Url          string
	Date         time.Time `gorm:"type:datetime"`
	TournamentId int       `                            selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-super > div:nth-child(1) > a"                                                        source:"attr=href" parser:"idParser"`
	Stage        Stage     `                            selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-super > div:nth-child(1) > a > div > div.match-header-event-series"                                    parser:"stageParser"`
	Team1Id      int       `gorm:"column:team_1_id"     selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1"                                        source:"attr=href" parser:"idParser"`
	Team2Id      int       `gorm:"column:team_2_id"     selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2"                                        source:"attr=href" parser:"idParser"`
	Team1Score   int       `gorm:"column:team_1_score"  selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span:nth-child(1)"`
	Team2Score   int       `gorm:"column:team_2_score"  selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > div > div.match-header-vs-score > div:nth-child(1) > span:nth-child(3)"`
	Team1Rating  int       `gorm:"column:team_1_rating" selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-1 > div > div.match-header-link-name-elo"                    parser:"ratingParser"`
	Team2Rating  int       `gorm:"column:team_2_rating" selector:"#wrapper > div.col-container > div.col.mod-3 > div.wf-card.match-header > div.match-header-vs > a.match-header-link.wf-link-hover.mod-2 > div > div.match-header-link-name-elo"                    parser:"ratingParser"`
}

type BanPickLogSchema struct {
	MatchId    int
	TeamId     *int
	VetoOrder  int
	MapId      int
	VetoAction VetoAction
}

type MatchMapSchema struct {
	MatchId       int
	MapId         int  `selector:"div.vm-stats-game-header > div.map > div:nth-child(1) > span"                       parser:"mapIdParser"`
	Duration      *int `selector:"div.vm-stats-game-header > div.map > div.map-duration"                              parser:"durationParser"`
	Team1Id       int  `                                                                                                                      gorm:"column:team_1_id"`
	Team2Id       int  `                                                                                                                      gorm:"column:team_2_id"`
	Team1DefScore int  `selector:"div.vm-stats-game-header > div:nth-child(1) > div:nth-child(2) > span.mod-ct"                               gorm:"column:team_1_def_score"`
	Team1AtkScore int  `selector:"div.vm-stats-game-header > div:nth-child(1) > div:nth-child(2) > span.mod-t"                                gorm:"column:team_1_atk_score"`
	Team1OTScore  int  `selector:"div.vm-stats-game-header > div:nth-child(1) > div:nth-child(2) > span.mod-ot"                               gorm:"column:team_1_ot_score"`
	Team2DefScore int  `selector:"div.vm-stats-game-header > div.team.mod-right > div:nth-child(1) > span.mod-ct"                             gorm:"column:team_2_def_score"`
	Team2AtkScore int  `selector:"div.vm-stats-game-header > div.team.mod-right > div:nth-child(1) > span.mod-t"                              gorm:"column:team_2_atk_score"`
	Team2OTScore  int  `selector:"div.vm-stats-game-header > div.team.mod-right > div:nth-child(1) > span.mod-ot"                             gorm:"column:team_2_ot_score"`
	TeamDefFirst  int  `selector:"div.vm-stats-game-header > div:nth-child(1) > div:nth-child(2) > span:nth-child(2)" parser:"defFirstParser"                                source:"attr=class"`
	TeamPick      *int `selector:"div.vm-stats-game-header > div.map > div:nth-child(1) > span > span.picked"         parser:"teamPickParser"                                source:"attr=class"`
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

type DuelKills struct {
	Team1PlayerKillsVsTeam2Player int `selector:"div:nth-child(1)" parser:"duelParser" gorm:"column:team_1_player_kills_vs_team_2_player"`
	Team2PlayerKillsVsTeam1Player int `selector:"div:nth-child(2)" parser:"duelParser" gorm:"column:team_2_player_kills_vs_team_1_player"`
}

type DuelFirstKills struct {
	Team1PlayerFirstKillsVsTeam2Player int `selector:"div:nth-child(1)" parser:"duelParser" gorm:"column:team_1_player_first_kills_vs_team_2_player"`
	Team2PlayerFirstKillsVsTeam1Player int `selector:"div:nth-child(2)" parser:"duelParser" gorm:"column:team_2_player_first_kills_vs_team_1_player"`
}

type DuelOpKills struct {
	Team1PlayerOpKillsVsTeam2Player int `selector:"div:nth-child(1)" parser:"duelParser" gorm:"column:team_1_player_op_kills_vs_team_2_player"`
	Team2PlayerOpKillsVsTeam1Player int `selector:"div:nth-child(2)" parser:"duelParser" gorm:"column:team_2_player_op_kills_vs_team_1_player"`
}

type PlayerDuelStatSchema struct {
	MatchId       int `gorm:"column:match_id"`
	MapId         int `gorm:"column:map_id"`
	Team1PlayerId int `gorm:"column:team_1_player_id"`
	Team2PlayerId int `gorm:"column:team_2_player_id"`
	DuelKills
	DuelFirstKills
	DuelOpKills
}
type PlayerHighlightSchema struct {
	MatchId         int
	MapId           int
	RoundNo         int
	TeamId          int
	PlayerId        int
	HighlightType   HighlightType
	PlayerAgainstId int
}

type RoundOverviewSchema struct {
	RoundNo   int       `selector:"div.rnd-num"`
	TeamWon   int       `selector:"div.rnd-sq.mod-win:nth-child(2)" source:"attr=class" parser:"teamWonParser"`
	TeamDef   int       `selector:"div.rnd-sq.mod-win"              source:"attr=class" parser:"teamDefParser"`
	WonMethod WonMethod `selector:"div.rnd-sq.mod-win > img"        source:"attr=src"   parser:"wonMethodParser"`
}

type RoundEconomySchema struct {
	Team1BuyType BuyType `selector:"div.rnd-sq:nth-child(3)" parser:"buyTypeParser" gorm:"column:team_1_buy_type"`
	Team2BuyType BuyType `selector:"div.rnd-sq:nth-child(4)" parser:"buyTypeParser" gorm:"column:team_2_buy_type"`
	Team1Bank    int     `selector:"div.bank:nth-child(2)"   parser:"balanceParser" gorm:"column:team_1_bank"`
	Team2Bank    int     `selector:"div.bank:nth-child(5)"   parser:"balanceParser" gorm:"column:team_2_bank"`
}

type RoundStatSchema struct {
	MatchId int
	MapId   int
	Team1Id int `gorm:"column:team_1_id"`
	Team2Id int `gorm:"column:team_2_id"`
	RoundOverviewSchema
	RoundEconomySchema
}

type TeamSchema struct {
	Id            int
	Name          string  `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.team-header > div.team-header-desc > div > div.team-header-name > h1"`
	ShorthandName *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.team-header > div.team-header-desc > div > div.team-header-name > h2"`
	Url           string
	ImgUrl        *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.team-header > div.wf-avatar.team-header-logo > div > img"             source:"attr=src"`
	CountryId     *int
	RegionId      *int
}

type TournamentSchema struct {
	Id        int
	Name      string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-event.mod-header.mod-full > div.event-header > div.event-desc > div > h1"`
	Url       string
	PrizePool int  `selector:"#wrapper > div.col-container > div > div.wf-card.mod-event.mod-header.mod-full > div.event-header > div.event-desc > div > div.event-desc-items > div:nth-child(2) > div.event-desc-item-value" parser:"moneyParser"`
	Tier1     bool `                                                                                                                                                                                                                               gorm:"column:tier_1"`
}
