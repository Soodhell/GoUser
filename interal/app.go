package interal

import (
	"User/interal/config"
	"User/interal/controllers"
	"User/interal/repositories"
	"User/interal/services"
	"User/pkg/db"
	"User/pkg/http"

	"github.com/gorilla/mux"
)

func StartApp() {

	var DB db.Postgres

	conf := config.GetConfig()
	err := DB.ConnectAndTest(conf.User, conf.Password, conf.Dbname, conf.Host, conf.Port)
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
