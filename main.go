package main

import (
	"github.com/cmiceli/configserver/lib"
	"log"
	"net/http"
)

func main() {
	store := configserver.NewFSStorage("./")
	httpServer := configserver.NewHTTPConfigServer(store)
	log.Fatal(http.ListenAndServe(":8000", httpServer))
}
