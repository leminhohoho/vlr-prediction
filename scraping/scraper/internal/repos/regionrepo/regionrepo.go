package regionrepo

import (
	"errors"
	"strings"

	"github.com/leminhohoho/vlr-prediction/scraping/scraper/internal/models"
	"gorm.io/gorm"
)

var ErrRowExist = errors.New("Row already exists")

type RegionRepo struct {
	tx *gorm.DB
}

func NewRegionRepo(tx *gorm.DB) *RegionRepo {
	return &RegionRepo{
		tx: tx,
	}
}

func (r *RegionRepo) InsertRegion(regionName string) (regionInfo models.RegionSchema, err error) {
	regionInfo.Name = regionName

	if rs := r.tx.Table("regions").Create(&regionInfo); rs.Error != nil {
		if strings.Contains(rs.Error.Error(), "UNIQUE constraint failed") {
			err = ErrRowExist
			return
		}

		err = rs.Error
	}

	return
}

func (r *RegionRepo) GetRegionByName(regionName string) (regionInfo models.RegionSchema, err error) {
	if rs := r.tx.Table("regions").Where("name = ?", regionName).First(&regionInfo); rs.Error != nil {
		err = rs.Error
	}

	return
}

func (r *RegionRepo) GetRegionById(regionId int) (regionInfo models.RegionSchema, err error) {
	if rs := r.tx.Table("regions").Where("id = ?", regionId).First(&regionInfo); rs.Error != nil {
		err = rs.Error
	}

	return
}
