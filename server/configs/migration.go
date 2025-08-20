package configs

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(localEnv *LocalEnv) *sql.DB {
	dbUri := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		localEnv.DatabaseUser,
		localEnv.DatabasePassword,
		localEnv.DatabaseHost,
		localEnv.DatabasePort,
		localEnv.DatabaseName,
	)
	db, err := sql.Open("postgres", dbUri)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)

	if err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migration: %v", err)
	}

	log.Println("Database migrations run successfully")

	return db
}
