//go:build !production

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/u104rak1/pocgo/internal/config"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/model"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/seed"
	"github.com/uptrace/bun"
)

const (
	schemaPath     = "internal/infrastructure/postgres/schema.sql"
	migrationsPath = "internal/infrastructure/postgres/migrations"
)

func main() {
	dsn := config.CreateDSN()
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	cmd := os.Args[1] + " " + os.Args[2]
	switch cmd {
	case "insert seed":
		insertSeedData()
	case "migrate up":
		migrateUp(m)
	case "migrate down":
		migrateDown(m)
	case "migrate reset":
		migrateReset(m)
	case "drop tables":
		dropTables(m)
	case "migrate refresh":
		updateSchemaAndGenerateMigrations(dsn, m)
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func insertSeedData() {
	db, err := config.LoadDB()
	if err != nil {
		log.Fatal(err)
	}
	seed.InsertMasterData(db)
	seed.InsertSeedData(db)
}

func migrateUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to migrate up: %v", err)
	}
	fmt.Println("Migrations applied successfully")
}

func migrateDown(m *migrate.Migrate) {
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to migrate down one step: %v", err)
	}
	fmt.Println("One step migration rolled back successfully")
}

func migrateReset(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to migrate down all: %v", err)
	}
	fmt.Println("All migrations rolled back successfully")
}

func dropTables(m *migrate.Migrate) {
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

	fmt.Println("Drop tables successfully")
}

func updateSchemaAndGenerateMigrations(dsn string, m *migrate.Migrate) {
	if err := checkAtlasInstalled(); err != nil {
		log.Fatal(err)
	}

	updateSchema()

	dropTables(m)

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

func foreignKeysToSQL() []byte {
	var data []byte
	for _, fk := range model.ForeignKeys {
		query := fmt.Sprintf(`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)`,
			fk.Table, fk.ConstraintName, fk.Column, fk.ReferencedTable, fk.ReferencedColumn)

		data = append(data, query...)
		data = append(data, ";\n"...)
	}
	return data
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
