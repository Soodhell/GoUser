package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func (p *Postgres) Connect() error {
	conn := "user=postgres password=132457689090iop dbname=school_go host=localhost port=5432 sslmode=disable"

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

func (p *Postgres) ConnectAndTest() error {
	err := p.Connect()
	if err != nil {
		return err
	}

	err = p.Ping()
	if err != nil {
		return err
	}

	return nil
}
