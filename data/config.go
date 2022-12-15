package data

import (
	"fmt"
	"main/data/artists"
	"main/data/labels"
	"main/data/platforms"
	"time"

	consulClient "github.com/punk-link/consul-client"
	"github.com/punk-link/logger"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func New(logger logger.Logger, consul consulClient.ConsulClient /*, appSecrets map[string]any*/) *gorm.DB {
	connectionString, err := getConnectionString(consul /*, appSecrets*/)
	db, err := openConnection(err, connectionString, logger)
	err = configureConnection(err, db)
	if err != nil {
		logger.LogFatal(err, err.Error())
	}

	db.Use(otelgorm.NewPlugin())

	autoMigrate(logger, db)

	return db
}

func autoMigrate(logger logger.Logger, db *gorm.DB) {
	err := migrateInternal(logger, nil, db, &labels.Label{}, &labels.Manager{})
	err = migrateInternal(logger, err, db, &artists.Artist{}, &artists.Release{}, &artists.ArtistReleaseRelation{})
	_ = migrateInternal(logger, err, db, &platforms.PlatformReleaseUrl{})
}

func configureConnection(err error, db *gorm.DB) error {
	if err != nil {
		return err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("postgres initialization failed: %v", err.Error())
	}

	sqlDb.SetConnMaxIdleTime(10)
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetConnMaxLifetime(time.Minute * 10)

	return err
}

func getConnectionString(consul consulClient.ConsulClient /*, appSecrets map[string]any*/) (string, error) {
	dbSettingsValue, err := consul.Get("DatabaseSettings")
	if err != nil {
		return "", fmt.Errorf("can't obtain Postgres configurations from Consul: %v", err.Error())
	}

	dbSettings := dbSettingsValue.(map[string]any)

	host := dbSettings["Host"].(string)
	port := dbSettings["Port"].(string)
	userName := dbSettings["UserName"].(string) //appSecrets["database-username"].(string)
	password := dbSettings["Password"].(string) //appSecrets["database-password"].(string)

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=punklink sslmode=disable TimeZone=UTC", host, port, userName, password), nil
}

func migrateInternal(logger logger.Logger, err error, db *gorm.DB, dst ...any) error {
	if err != nil {
		logger.LogFatal(err, err.Error())
	}

	return db.AutoMigrate(dst...)
}

func openConnection(err error, connectionString string, logger logger.Logger) (*gorm.DB, error) {
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connectionString,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger.New(logger, gormLogger.Config{
			LogLevel: gormLogger.Error,
		}),
	})

	if err != nil {
		err = fmt.Errorf("postgres initialization failed: %v", err.Error())
	}

	return db, err
}
