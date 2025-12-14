package interal

import (
	"User/interal/config"
	"User/interal/controllers"
	"User/interal/repositories"
	"User/interal/services"
	"User/pkg/db"
	"User/pkg/http"
)

func StartApp() {

	var DB db.Postgres
	err := DB.ConnectAndTest()
	defer DB.Close()

	if err != nil {
		panic(err)
	}

	rep := repositories.StartRepository(&DB)
	ser := services.StartService(*rep)
	con := controllers.StartController(*ser)

	var listController []config.Controller
	listController = append(listController, con)

	http.Run(config.Search(listController))

}
