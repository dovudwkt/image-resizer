package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/dovudwkt/playground/server/handlers"
	service "github.com/dovudwkt/playground/server/services"
)

func main() {

	var port = 3001
	var addr = fmt.Sprintf(":%d", port)

	// init service
	svc := service.New(service.Options{})

	// routes
	http.Handle("/images/resize", handler.ImageHTTPHandler{
		Service: svc,
	})
	http.Handle("/cutter", handler.ImageFromURLHTTPHandler{
		Service: svc,
	})

	log.Print("Starting server at port ", port)
	http.ListenAndServe(addr, nil)

}
