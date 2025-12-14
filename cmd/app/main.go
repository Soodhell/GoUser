package main

import (
	"User/interal"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "User/docs"
)

// @title User API
// @version 1.0
// @description API server

// @host localhost:12543
func main() {

	go interal.StartApp()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println("Приложение закончило работу")

}
