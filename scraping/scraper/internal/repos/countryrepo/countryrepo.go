package countryrepo

import (
	"errors"
	"strings"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/gorm"
)

var ErrRowExists = errors.New("Row already exists")

type CountryRepo struct {
	tx *gorm.DB
}

func NewCountryRepo(tx *gorm.DB) *CountryRepo {
	return &CountryRepo{
		tx: tx,
	}
}

func (c *CountryRepo) InsertCountry(
	regionId int,
	countryName string,
) (countryInfo models.CountrySchema, err error) {
	countryInfo.Name = countryName
	countryInfo.RegionId = regionId

	if rs := c.tx.Table("countries").Create(&countryInfo); rs.Error != nil {
		if strings.Contains(rs.Error.Error(), "UNIQUE constraint failed") {
			err = ErrRowExists
			return
		}

		err = rs.Error
	}

	return
}

func (c *CountryRepo) GetCountryByName(
	countryName string,
) (countryInfo models.CountrySchema, err error) {
	if rs := c.tx.Table("countries").Where("name = ?", countryName).First(&countryInfo); rs.Error != nil {
		err = rs.Error
	}

	return
}

func (c *CountryRepo) GetCountryById(
	countryId int,
) (countryInfo models.CountrySchema, err error) {
	if rs := c.tx.Table("countries").Where("id= ?", countryId).First(&countryInfo); rs.Error != nil {
		err = rs.Error
	}

	return
}
