package database

import (
	"fmt"
	"time"

	env "github.com/andre-felipe-wonsik-alves/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	host := env.GetEnv("DB_HOST", "localhost")
	port := env.GetEnv("DB_PORT", "5432")
	user := env.GetEnv("DB_USER", "app_user")
	password := env.GetEnv("DB_PASSWORD", "app_password")
	databaseName := env.GetEnv("DB_NAME", "app_db")
	dev := env.GetEnv("DEVELOPMENT", "true")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, databaseName,
	)

	logLevel := logger.Warn

	if dev == "true" {
		logLevel = logger.Info
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
