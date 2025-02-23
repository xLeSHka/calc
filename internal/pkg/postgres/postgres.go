package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(config config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB, config.PostgresHost, config.PostgresPort),
		PreferSimpleProtocol: true,
	}), gormConfig)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	DB, err := db.DB()
	if err != nil {
		return err
	}
	if err := goose.Up(DB, "migrations"); err != nil {
		return err
	}
	return nil
}
