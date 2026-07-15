package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syncova-todo/config"
	domain "syncova-todo/domain/base"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresConnection(cfg *config.Config, withLog bool) (*PostgresDB, error) {
	dsn := "postgresql://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=require&channel_binding=require"

	logLevel := logger.Silent
	if cfg.ENV == "development" && withLog {
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

func (p *PostgresDB) Paginated(ctx context.Context, request *domain.BasePaginationRequest, models interface{}, searchColumns []string) (*domain.BaseListResponse, error) {
	page := request.Page
	if page <= 0 {
		page = 1
	}

	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}

	query := p.DB.WithContext(ctx).Model(models)
	if request.Search != "" && len(searchColumns) > 0 {
		search := "%" + strings.TrimSpace(request.Search) + "%"
		conditions := make([]string, 0, len(searchColumns))
		args := make([]interface{}, 0, len(searchColumns))

		for _, column := range searchColumns {
			if strings.TrimSpace(column) == "" {
				continue
			}

			conditions = append(conditions, column+" ILIKE ?")
			args = append(args, search)
		}

		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "), args...)
		}
	}

	sortBy := request.SortBy
	if strings.TrimSpace(sortBy) == "" {
		sortBy = "id"
	}

	result := query.
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: sortBy},
			Desc:   strings.ToLower(request.Sort) == "desc",
		}).
		Limit(perPage).
		Offset((page - 1) * perPage).
		Find(models)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	return &domain.BaseListResponse{
		Data:    models,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}
