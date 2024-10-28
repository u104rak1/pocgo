package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func LoadDB() (*bun.DB, error) {
	db, err := newDBConnection()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}

func CreateDSN() string {
	env := NewEnv()
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Tokyo",
		env.POSTGRES_USER, env.POSTGRES_PASSWORD, env.POSTGRES_HOST, env.POSTGRES_PORT, env.POSTGRES_DBNAME, env.POSTGRES_SSLMODE)
}

func newDBConnection() (*bun.DB, error) {
	dsn := CreateDSN()
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)

	if err := sqldb.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	// Register many-to-many intermediate tables here
	// db.RegisterModel((*model.ManyToMany)(nil))

	log.Println("Successfully connected to database")
	return db, nil
}

func CloseDB(db *bun.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
	log.Println("Successfully closed database")
}
