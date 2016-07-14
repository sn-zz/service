// sn - https://github.com/sn
package main

import (
    "net/http"

    "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)

    router.Handle("/", Index).Methods(http.MethodGet)
    router.Handle("/auth", Auth).Methods(http.MethodPost)

    router.Handle("/users", UserIndex).Methods(http.MethodGet)
    router.Handle("/users/{userId}", UserShow).Methods(http.MethodGet)
    router.Handle("/users/{userId}", UserCreate).Methods(http.MethodPost)
    router.Handle("/users/{userId}", UserUpdate).Methods(http.MethodPut)
    router.Handle("/users/{userId}", UserPatch).Methods(http.MethodPatch)
    router.Handle("/users/{userId}", UserDelete).Methods(http.MethodDelete)

    return router
}
