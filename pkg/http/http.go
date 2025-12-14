package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func Run(mux *mux.Router) {

	server := http.Server{
		Addr:           ":12543",
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("http server start, port: 12543")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
