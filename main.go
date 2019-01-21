package main

import (
	"github.com/cmiceli/configserver/lib"
	"log"
	"net/http"
)

func main() {
	store := configserver.NewMemStorage()
	httpServer := configserver.NewHTTPConfigServer(store)
	log.Fatal(http.ListenAndServe(":8000", httpServer))
}
