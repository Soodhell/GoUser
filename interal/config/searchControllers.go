package config

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	_ "User/docs"
)

type Controller interface {
	GetController() map[string]func(http.ResponseWriter, *http.Request)
	GetMethod() map[string]string
}

func Search(listController []Controller) *mux.Router {

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	for _, controller := range listController {
		for controllerName, controllerMethod := range controller.GetController() {
			router.HandleFunc(controllerName, controllerMethod).Methods(controller.GetMethod()[controllerName])
		}
	}

	return router
}
