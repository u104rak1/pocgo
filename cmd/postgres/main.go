package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ucho456job/pocgo/internal/config"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/seed"
	"github.com/uptrace/bun"
)

const (
	schemaPath     = "internal/infrastructure/postgres/schema.sql"
	migrationsPath = "internal/infrastructure/postgres/migrations"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected 'generate', 'migrate up', 'migrate downall', 'migrate downone', 'migrate generate' 'migrate reset' or 'seed' subcommands")
	}

	switch strings.ToLower(os.Args[1]) {
	case "generate":
		updateSchema()
	case "migrate":
		if len(os.Args) < 3 {
			log.Fatal("Expected 'up', 'downall', 'downone', or 'reset' subcommands for 'migrate'")
		}
		switchMigrateCommand(os.Args[2])
	case "seed":
		InsertSeedData()
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
	data = append(data, foreignKeysToSQL()...)

	os.WriteFile(schemaPath, data, 0777)
	fmt.Println("Successfully updated schema.sql")
}

func modelsToByte(db *bun.DB, models []interface{}) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model)
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

// .WithForeignKeys() has oneのリレーションがうまくいかないので、外部キー制約は手動で追加
func foreignKeysToSQL() []byte {
	var data []byte
	for _, fk := range model.ForeignKeys {
		query := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s) ON DELETE CASCADE`,
			fk.Table, fk.ConstraintName, fk.Column, fk.ReferencedTable, fk.ReferencedColumn)

		data = append(data, query...)
		data = append(data, ";\n"...)
	}
	return data
}

func switchMigrateCommand(action string) {
	dsn := config.CreateDSN()
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
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
	case "generate":
		generateMigrations(dsn, m)
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

func generateMigrations(dsn string, m *migrate.Migrate) {
	if err := checkAtlasInstalled(); err != nil {
		log.Fatal(err)
	}

	updateSchema()

	resetDatabase(m)

	atlasCmd := []string{
		"migrate", "diff", "migration",
		"--dir", fmt.Sprintf("file://%s?format=golang-migrate", migrationsPath),
		"--to", fmt.Sprintf("file://%s", schemaPath),
		"--dev-url", dsn,
	}
	if err := runCommand("atlas", atlasCmd...); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migration file generated successfully")
}

func checkAtlasInstalled() error {
	_, err := exec.LookPath("atlas")
	if err != nil {
		return fmt.Errorf("atlas is not installed. Please install it by running:\ncurl -sSf https://atlasgo.sh | sh")
	}
	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	return cmd.Run()
}

func InsertSeedData() {
	db, err := config.LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	seed.InsertMasterData(db)
	seed.InsertSeedData(db)
}
