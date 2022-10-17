package db

import (
	"github.com/recative/recative-backend-sdk/pkg/config"
	"github.com/recative/recative-backend-sdk/pkg/logger"
	"go.uber.org/zap"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type AutoMigrater interface {
	AutoMigrate()
}

type Config struct {
	IsAutoMigrate bool   `env:"IS_AUTO_MIGRATE" envDefault:"true"`
	PgsqlUri      string `env:"PGSQL_URI"`
	PgsqlMaxIdle  int    `env:"PGSQL_MAX_IDLE" envDefault:"50"`
	PgsqlMaxOpen  int    `env:"PGSQL_MAX_OPEN" envDefault:"50"`
}

func New(config_ Config) *gorm.DB {
	var gormLogger gormlogger.Interface

	if config.Environment() == config.Prod {
		gormLogger = NewProductionGormLoggerConfig().BuildWith(logger.RawLogger().Sugar())
	} else {
		gormLogger = NewDevelopmentGormLoggerConfig().BuildWith(logger.RawLogger().Sugar())
	}

	db, err := gorm.Open(driver.Open(config_.PgsqlUri), &gorm.Config{Logger: gormLogger})
	if err != nil {
		panic("Open PostgreSQL DB failed: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic("init database failed", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(config_.PgsqlMaxIdle)
	sqlDB.SetMaxOpenConns(config_.PgsqlMaxOpen)

	return db
}
