module github.com/leminhohoho/vlr-prediction/scraping/scraper

go 1.24.2

require (
	github.com/PuerkitoBio/goquery v1.10.3
	github.com/gocolly/colly v1.2.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/joho/godotenv v1.5.1
	github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx v0.0.0
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/sirupsen/logrus v1.9.3
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/antchfx/htmlquery v1.3.4 // indirect
	github.com/antchfx/xmlquery v1.4.4 // indirect
	github.com/antchfx/xpath v1.3.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/leminhohoho/vlr-prediction/scraping/pkgs/htmlx v0.0.0 => ../pkgs/htmlx
