package data

import (
	"fmt"
	"main/data/artists"
	"main/data/labels"
	"main/infrastructure/consul"
	"main/services/common"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigureDatabase(logger *common.Logger, consul *consul.ConsulClient) {
	dbSettings := consul.Get("DatabaseSettings").(map[string]interface{})

	host := dbSettings["Host"].(string)
	port := dbSettings["Port"].(string)
	userName := dbSettings["Username"].(string)
	password := dbSettings["Password"].(string)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=punklink sslmode=disable TimeZone=UTC", host, port, userName, password),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
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

func AutoMigrate(logger *common.Logger) {
	err := DB.AutoMigrate(&labels.Label{}, &labels.Manager{})
	if err != nil {
		logger.LogFatal(err, err.Error())
	}

	err = DB.AutoMigrate(&artists.Artist{}, &artists.Release{})
	if err != nil {
		logger.LogFatal(err, err.Error())
	}
}

var DB *gorm.DB
