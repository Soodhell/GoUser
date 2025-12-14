package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func (p *Postgres) Connect(user string, password string, dbname string, host string, port string) error {
	//conn := "user=postgres password=132457689090iop dbname=school_go host=localhost port=5432 sslmode=disable"
	conn := "user=" + user + " password=" + password + " dbname=" + dbname + " host=" + host + " port=" + port + " sslmode=disable"

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
	err := p.Connect(user, password, dbname, host, port)
	if err != nil {
		return err
	}

	err = p.Ping()
	if err != nil {
		return err
	}

	return nil
}
