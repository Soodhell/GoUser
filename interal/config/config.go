package config

import "flag"

type Config struct {
	User     string
	Password string
	Dbname   string
	Host     string
	Port     string
}

func GetConfig() Config {

	var config Config

	flag.StringVar(&config.User, "user", "", "postgres user")
	flag.StringVar(&config.Password, "password", "", "postgres password")
	flag.StringVar(&config.Dbname, "dbname", "", "postgres database name")
	flag.StringVar(&config.Host, "host", "", "postgres host")
	flag.StringVar(&config.Port, "port", "", "postgres port")

	flag.Parse()

	return config
}
