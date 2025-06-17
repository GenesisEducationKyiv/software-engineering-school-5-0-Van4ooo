package db

import (
	"database/sql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const migrationsDir = "migrations"

func Init(postgres config.Postgres) {
	dbConn := openGormDB(postgres.URL)
	DB = dbConn
	sqlDB := getSQLDB(dbConn)
	runMigrations(sqlDB)
}

func openGormDB(databaseURL string) *gorm.DB {
	dbConn, err := gorm.Open(gormDriver.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect via GORM: %v", err)
	}
	return dbConn
}

func getSQLDB(dbConn *gorm.DB) *sql.DB {
	sqlDB, err := dbConn.DB()
	if err != nil {
		log.Fatalf("failed to get generic DB: %v", err)
	}
	return sqlDB
}

func runMigrations(sqlDB *sql.DB) {
	driver, err := migrateDriver.WithInstance(sqlDB, &migrateDriver.Config{})
	if err != nil {
		log.Fatalf("migration driver init error: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir, "postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration up error: %v", err)
	}
}
