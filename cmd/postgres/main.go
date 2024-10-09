package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ucho456job/pocgo/internal/config"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected 'generate, 'migrate up', 'migrate downall', 'migrate downone', or 'migrate reset' subcommands")
	}

	switch strings.ToLower(os.Args[1]) {
	case "generate":
		updateSchema()
	case "migrate":
		if len(os.Args) < 3 {
			log.Fatal("Expected 'up', 'downall', 'downone', or 'reset' subcommands for 'migrate'")
		}
		migrateCommand(os.Args[2])
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func updateSchema() {
	db, err := config.LoadDB()
	if err != nil {
		log.Fatal(err)
	}

	var data []byte
	data = append(data, modelsToByte(db, model.Models)...)
	data = append(data, indexesToByte(db, model.AllIdxCreators())...)

	os.WriteFile("internal/infrastructure/postgres/schema.sql", data, 0777)
	fmt.Println("Successfully updated schema.sql")
}

func modelsToByte(db *bun.DB, models []interface{}) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys()
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

func indexesToByte(db *bun.DB, idxCreators []model.IndexQueryCreators) []byte {
	var data []byte
	for _, idxCreator := range idxCreators {
		idx := idxCreator(db)
		rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

func migrateCommand(action string) {
	dsn := config.CreateDSN()
	m, err := migrate.New("file://internal/infrastructure/postgres/migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch strings.ToLower(action) {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate up: %v", err)
		}
		fmt.Println("Migrations applied successfully")
	case "downall":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate down all: %v", err)
		}
		fmt.Println("All migrations rolled back successfully")
	case "downone":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate down one step: %v", err)
		}
		fmt.Println("One step migration rolled back successfully")
	case "reset":
		resetDatabase(m)
	default:
		log.Fatalf("Unknown migrate action: %s", action)
	}
}

func resetDatabase(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to migrate down all: %v", err)
	}

	db, err := config.LoadDB()
	if err != nil {
		log.Fatalf("Failed to load database: %v", err)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS schema_migrations;"); err != nil {
		log.Fatalf("Failed to drop schema_migrations table: %v", err)
	}

	fmt.Println("Database reset successfully")
}
