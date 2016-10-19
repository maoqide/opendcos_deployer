package main

import (
	"deployer/services"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {

	wsContainer := restful.NewContainer()
	u := services.Resource{}
	u.Register(wsContainer)

	log.Printf("start listening on localhost:28080")
	server := &http.Server{Addr: ":28080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
