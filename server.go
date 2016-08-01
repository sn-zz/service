// Package server is the main program that initiates the server.
//
// sn - https://github.com/sn
package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/sn/service/router"
)

func main() {
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router)))
}
