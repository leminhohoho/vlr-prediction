package models

type CountrySchema struct {
	Id       int    `gorm:"column:id;primaryKey;autoIcrement"`
	Name     string `gorm:"column:name"`
	RegionId int    `gorm:"column:region_id"`
}

type RegionSchema struct {
	Id   int    `gorm:"column:id;primaryKey;autoIcrement"`
	Name string `gorm:"column:name"`
}

type PlayerSchema struct {
	Id        int
	Name      string  `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h1"`
	RealName  *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div:nth-child(1) > h2"`
	Url       string
	ImgUrl    *string `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div.wf-avatar.mod-player > div > img"     source:"attr=src"`
	CountryId *int    `selector:"#wrapper > div.col-container > div > div.wf-card.mod-header.mod-full > div.player-header > div:nth-child(2) > div.ge-text-light"                       parser:"countryIdParser"`
}
