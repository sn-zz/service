// sn - https://github.com/sn
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
)

// GET /index
var Index = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if auth := r.Header["Authorization"]; auth != nil {
        if session := FindSession(auth[0]); session.Id != "" {
            u := FindUserById(session.UserId)
            fmt.Fprintf(w, "Welcome, %s!\n", u.Username)
            UpdateSessionTime(session.Id)
            return
        }
    }
    fmt.Fprint(w, "Welcome!\n")
})

// POST /auth
var Auth = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var user User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.Unmarshal(body, &user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    refUser := FindUserById(user.Id)
    if len(refUser.Id) > 0 {
        if CheckPassword(refUser, user.Password) {
            w.WriteHeader(http.StatusOK)
            s, _ := CreateSession(refUser.Id)
            fmt.Fprintf(w, "%s", GenerateSha1Hash(string(s.Id)))
            return
        }
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    fmt.Printf("userId %s, password %s", user.Id, user.Password)

    w.WriteHeader(http.StatusNotFound)
})

// GET /users
var UserIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(users); err != nil {
        panic(err)
    }
})

// GET /users/:userId
var UserShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := uuid(vars["userId"])
    user := FindUserById(userId)
    if len(user.Id) > 0 {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(user); err != nil {
            panic(err)
        }
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not Found")
})

// POST /users/:userId
var UserCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var user User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    if user := IsUsernameTaken(user.Username); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Username is taken")
        return
    }
    if user := IsAddressTaken(user.Address.Address); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Address is taken")
        return
    }

    if user, err := CreateUser(user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(user); err != nil {
            fmt.Fprint(w, err)
            panic(err)
        }
        return
    }
})

// PUT /users/:userId
var UserUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := uuid(vars["userId"])

    var user User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
    if user := IsUsernameTaken(user.Username); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Username is taken")
        return
    }
    if user := IsAddressTaken(user.Address.Address); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Address is taken")
        return
    }

    if user, err := UpdateUser(userId, user); err == nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(user); err != nil {
            panic(err)
        }
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not Found")
})

// PATCH /users/:userId
var UserPatch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := uuid(vars["userId"])

    var user User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, "Bad request")
        return
    }
    if user := IsUsernameTaken(user.Username); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Username is taken")
        return
    }
    if user := IsAddressTaken(user.Address.Address); user {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusConflict)
        fmt.Fprint(w, "Address is taken")
        return
    }

    if user, err := PatchUser(userId, user); err == nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(user); err != nil {
            panic(err)
        }
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not found")
})

// DELETE /users/:userId
var UserDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := uuid(vars["userId"])

    if err := DeleteUser(userId); err == nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusNoContent)
        fmt.Fprint(w, err)
        return
    }

    // If we didn't find it, 404
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not found")
})
