package piper

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "embed"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var (
	ErrIncorrectSchema = errors.New("Missmatch database schema")
)

//go:embed cache_storage_schema.sql
var piperCacheDbSchema string

var gormConfig = gorm.Config{
	Logger: logger.Default.LogMode(logger.Silent), // Disable all logs
}

func init() {
	piperCacheDbSchema = strings.Join(strings.Fields(piperCacheDbSchema), " ")
}

// Cache is the interface implemented by objects those can perform actions of a key:value cache database.
type Cache interface {
	// Set take a key name, a value of []byte type and an expiration duration and saved it.
	// The value is later retrievable by using [Cache.Get] with the key name.
	// The method return value if the saving process encounter errors.
	Set(string, []byte, time.Duration) error
	// Get take a key name and return the corresponding value to that key.
	// It return a 2nd value as an error if it failed to get the value from the cache storage.
	Get(string) ([]byte, error)
	// Delete take a key name and remove the value associated with the key name.
	// It return 2nd value as error if it failed to remove the value and the corrsponding key from the cache storage.
	// If the value didn't exist, the cache storage will be unchanged, and no error will be returned.
	Delete(string) error
}

// PiperCache is an implementation of [Cache], used as default cache database.
// It is safe for concurrent usage.
type PiperCache struct {
	mu sync.Mutex

	src string
	db  *gorm.DB
}

type piperCacheTable struct {
	Key     string `gorm:"column:KEY"`
	Value   []byte `gorm:"column:value"`
	ExpDate *int64 `gorm:"column:expiration_timestamp"`
}

// NewCacheDb return an uninitialized cache storage.
// It take a path to .db file as a SQLite backend.
// If the file does not exists, it will attempt to create a new one.
// If exists the file stay unchanged. To manually setup a cache storage ready to be used, use [PiperCache.Setup].
// To validate whether the storage is ready to be used, use [PiperCache.Validate].
func NewCacheDb(src string) (*PiperCache, error) {
	db, err := gorm.Open(sqlite.Open(src), &gormConfig)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(src)
	if err != nil {
		return nil, err
	}

	return &PiperCache{db: db, src: absPath}, nil
}

// Validate check whether the cache storage is ready to be used.
// It return [ErrIncorrectSchema] if the SQLite database schema does not match or fetching the schema return error.
func (c *PiperCache) Validate() error {
	cmd := exec.Command("sqlite3", c.src, ".schema")
	schemaBytes, err := cmd.Output()
	if err != nil {
		return err
	}

	schema := strings.Join(strings.Fields(string(schemaBytes)), " ")

	if strings.TrimSpace(piperCacheDbSchema) != strings.TrimSpace(schema) {
		return ErrIncorrectSchema
	}

	return nil
}

// Setup initialize the cache storage by delete the linked database then creating and setting up a new one.
// It return error if recreating the database or executing the schema on database return error.
func (c *PiperCache) Setup() error {
	var err error

	if err := os.Remove(c.src); err != nil {
		return err
	}

	c.db, err = gorm.Open(sqlite.Open(c.src), &gormConfig)
	if err != nil {
		return err
	}

	if err = c.db.Exec(piperCacheDbSchema).Error; err != nil {
		return err
	}

	return nil
}

func (c *PiperCache) update() error {
	var rows []piperCacheTable

	if err := c.db.Table("cache").Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		expDate := time.Unix(*row.ExpDate, 0)
		if expDate.After(time.Now()) {
			continue
		}

		if err := c.db.Table("cache").Where("KEY = ?", row.Key).Delete(&piperCacheTable{}).Error; err != nil {
			return err
		}
	}

	return nil
}

// Set implement the [Cache] interface
func (c *PiperCache) Set(key string, val []byte, duration time.Duration) error {
	if err := c.update(); err != nil {
		return err
	}

	row := piperCacheTable{Key: key, Value: val}
	var expDate int64

	if int(duration) != 0 {
		expDate = time.Now().Add(duration).Unix()
		row.ExpDate = &expDate
	}

	if err := c.db.Table("cache").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "KEY"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&row).Error; err != nil {
		return err
	}

	return nil
}

// Get implement the [Cache] interface
func (c *PiperCache) Get(key string) ([]byte, error) {
	if err := c.update(); err != nil {
		return nil, err
	}

	var rs piperCacheTable

	if err := c.db.Table("cache").Where("key = ?", key).First(&rs).Error; err != nil {
		return nil, err
	}

	return rs.Value, nil
}

// Delete implement the [Cache] interface
func (c *PiperCache) Delete(key string) error {
	if err := c.update(); err != nil {
		return err
	}

	if err := c.db.Table("cache").Where("key = ?").Delete(&piperCacheTable{}).Error; err != nil {
		return err
	}

	return nil
}
