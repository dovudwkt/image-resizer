package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/dovudwkt/image_resizer/handler"
	service "github.com/dovudwkt/image_resizer/service"
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

	log.Print("Starting server at port ", port)
	http.ListenAndServe(addr, nil)

}
