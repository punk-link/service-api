package data

import (
	"fmt"
	"main/data/labels"
	"main/utils"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigureDatabase() {
	host := utils.GetEnvironmentVariable("DB_HOST")
	port := utils.GetEnvironmentVariable("DB_PORT")
	userName := utils.GetEnvironmentVariable("DB_USERNAME")
	password := utils.GetEnvironmentVariable("DB_PASSWORD")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=punklink sslmode=disable TimeZone=UTC", host, port, userName, password),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msgf("Postgres initialization failed: %v", err.Error())
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msgf("Postgres initialization failed: %v", err.Error())
	}

	sqlDb.SetConnMaxIdleTime(10)
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)

	DB = db

	AutoMigrate()
}

func AutoMigrate() {
	err := DB.AutoMigrate(&labels.Label{}, &labels.Manager{})

	if err != nil {
		log.Fatal().AnErr(err.Error(), err)
	}
}

var DB *gorm.DB
