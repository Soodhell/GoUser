package interal

import (
	"User/interal/controllers"
	"User/interal/repositories"
	"User/interal/services"
	"User/pkg/db"
	"User/pkg/http"
	"flag"

	"github.com/gorilla/mux"
)

func StartApp() {

	var DB db.Postgres

	var user string
	var password string
	var dbname string
	var host string
	var port string

	flag.StringVar(&user, "user", "", "postgres user")
	flag.StringVar(&password, "password", "", "postgres password")
	flag.StringVar(&dbname, "dbname", "", "postgres database name")
	flag.StringVar(&host, "host", "", "postgres host")
	flag.StringVar(&port, "port", "", "postgres port")

	flag.Parse()

	err := DB.ConnectAndTest(user, password, dbname, host, port)
	defer DB.Close()

	if err != nil {
		panic(err)
	}

	rep := repositories.StartRepository(&DB)
	ser := services.StartService(*rep)
	con := controllers.StartController(*ser)

	router := mux.NewRouter()
	con.SettingRouter(router)

	http.Run(router)

}
