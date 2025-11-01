package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/goggle-source/grpc-servic/sso/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	cfg := config.MustLoad()

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.NameDB)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	var MigrationsTableName, migrationsPath string
	flag.StringVar(&MigrationsTableName, "migrations-table", "migrations", "set name for migrations table")
	flag.StringVar(&migrationsPath, "migrations-path", ".\\migrations", "get path to migrations")
	flag.Parse()
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: MigrationsTableName, // Ваше кастомное имя таблицы
	})
	if err != nil {
		panic(err)
	}

	m, _ := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	fmt.Println("Migrations applied successfully!")
}
