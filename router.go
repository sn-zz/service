// sn - https://github.com/sn
package main

import "github.com/gorilla/mux"

// NewRouter sets up the URL routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/", Index).Methods("GET")
	router.Handle("/auth", Auth).Methods("POST")

	router.Handle("/users", UserIndex).Methods("GET")
	router.Handle("/users", UserCreate).Methods("POST")
	router.Handle("/users/{userId}", UserShow).Methods("GET")
	router.Handle("/users/{userId}", UserUpdate).Methods("PUT")
	router.Handle("/users/{userId}", UserPatch).Methods("PATCH")
	router.Handle("/users/{userId}", UserDelete).Methods("DELETE")

	return router
}
