package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func (p *Postgres) Connect(conn string) error {

	var err error

	p.DB, err = sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	return nil

}

func (p *Postgres) Ping() error {
	err := p.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Close() {
	err := p.DB.Close()
	if err != nil {
		panic(err)
	}
}

func (p *Postgres) ConnectAndTest(user string, password string, dbname string, host string, port string) error {
	//conn := "user=postgres password=132457689090iop dbname=school_go host=localhost port=5432 sslmode=disable"
	//go run cmd/app/main.go --user=postgres --password=132457689090iop --dbname=school_go --host=localhost --port=5432
	conn := "user=" + user + " password=" + password + " dbname=" + dbname + " host=" + host + " port=" + port + " sslmode=disable"

	err := p.Connect(conn)
	if err != nil {
		return err
	}

	err = p.Ping()
	if err != nil {
		return err
	}

	println("Successfully connected to database")

	driver, err := postgres.WithInstance(p.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	println("Successfully migrated database schema")

	return nil
}
