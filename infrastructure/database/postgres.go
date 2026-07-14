package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"syncova-todo/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresConnection(cfg *config.Config) (*PostgresDB, error) {
	dsn := "postgresql://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=require&channel_binding=require"

	logLevel := logger.Silent
	if cfg.ENV == "development" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		PrepareStmt:            true,  // Enable prepared statement cache
		SkipDefaultTransaction: false, // Use transactions by default
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	sqlDB.SetMaxIdleConns(25)                 // Maximum number of idle connections in the pool
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: db}, nil

}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgresDB) Transaction(fn func(tx *gorm.DB) error) error {
	return p.DB.Transaction(fn)
}

func (p *PostgresDB) AutoMigrate(models ...interface{}) error {
	return p.DB.AutoMigrate(models...)
}

func (p *PostgresDB) WithContext(ctx context.Context) *gorm.DB {
	return p.DB.WithContext(ctx)
}
