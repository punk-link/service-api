package data

import (
	"fmt"
	"main/data/artists"
	"main/data/labels"
	"main/data/platforms"
	"time"

	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func ConfigureDatabase(logger logger.Logger, consul *consulClient.ConsulClient) {
	dbSettingsValue, err := consul.Get("DatabaseSettings")
	if err != nil {
		logger.LogError(err, "Postgres initialization failed: %v", err.Error())
	}

	dbSettings := dbSettingsValue.(map[string]interface{})

	host := dbSettings["Host"].(string)
	port := dbSettings["Port"].(string)
	userName := dbSettings["Username"].(string)
	password := dbSettings["Password"].(string)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=punklink sslmode=disable TimeZone=UTC", host, port, userName, password),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger.New(logger, gormLogger.Config{
			LogLevel: gormLogger.Error,
		}),
	})
	if err != nil {
		logger.LogError(err, "Postgres initialization failed: %v", err.Error())
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.LogError(err, "Postgres initialization failed: %v", err.Error())
	}

	sqlDb.SetConnMaxIdleTime(10)
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)

	DB = db

	AutoMigrate(logger)
}

func AutoMigrate(logger logger.Logger) {
	err := migrate(logger, nil, &labels.Label{}, &labels.Manager{})
	err = migrate(logger, err, &artists.Artist{}, &artists.Release{}, &artists.ArtistReleaseRelation{})
	_ = migrate(logger, err, &platforms.PlatformReleaseUrl{})
}

func migrate(logger logger.Logger, err error, dst ...interface{}) error {
	if err != nil {
		logger.LogFatal(err, err.Error())
		return err
	}

	return DB.AutoMigrate(dst...)
}

var DB *gorm.DB
