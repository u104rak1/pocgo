package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func ConnectionDB() *bun.DB {
	env := NewEnv()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env.MYSQL_USER, env.MYSQL_PASSWORD, env.MYSQL_HOST, env.MYSQL_PORT, env.MYSQL_DBNAME)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open MySQL connection: %v", err)
	}

	db := bun.NewDB(sqlDB, mysqldialect.New())

	return db
}
