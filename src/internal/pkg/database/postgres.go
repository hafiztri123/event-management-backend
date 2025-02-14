package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hafiztri123/src/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dataSource := setDataSource(cfg)
	logger := setLogger()
	db := openConnection(dataSource, logger)
	setConnectionPool(db)
	return db, nil
}

func setDataSource(cfg *config.DatabaseConfig) string {
	fmt.Println(cfg)
	return fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.Host,
        cfg.User,
        cfg.Password,
        cfg.DBName,
        cfg.Port,
        cfg.SSLMode,
    )
}

func setLogger() logger.Interface{
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful: true,
		},
	)
}

func openConnection(datasource string, logger logger.Interface) (*gorm.DB) {
	log.Printf("[DEBUG] Attempting to connect with datasource: %s\n", datasource)
	db, err := gorm.Open(postgres.Open(datasource), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		log.Fatal("[FAIL] attempt to connect database has failed: %w", err)
	}

	return db
}

func setConnectionPool(db *gorm.DB)  {
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal("[FAIL] failed to get database instance: %w", err)
	}

	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Hour)
}