// Package router contains endpoint information for the service.
//
// sn - https://github.com/sn
package router

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"time"

	"github.com/gorilla/mux"
	"github.com/sn/service/helpers"
	"github.com/sn/service/session"
	"github.com/sn/service/types"
	"github.com/sn/service/user"
)

// Index handles GET /index
var Index = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if auth := r.Header["Authorization"]; auth != nil {
		if s := session.Find(auth[0]); s.ID != "" {
			if time.Now().Before(s.Expires) {
				u := user.FindByID(s.UserID)
				fmt.Fprintf(w, "Welcome, %s!\n", u.Username)
				err := session.Bump(s.ID)
				if err != nil {
					panic(err)
				}
				return
			} else {
				session.Expire(s.ID)
			}
		}
	}
	fmt.Fprint(w, "Welcome!\n")
})

// Auth handles POST /auth
var Auth = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	u := user.User{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.Unmarshal(body, &u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	refUser := user.FindByID(u.ID)
	if len(refUser.ID) > 0 {
		if user.CheckPassword(refUser, u.Password) {
			w.WriteHeader(http.StatusOK)
			s := session.Create(refUser.ID)
			fmt.Fprintf(w, "%s", helpers.GenerateSha1Hash(string(s.ID)))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNotFound)
})

// UserIndex handles GET /users
var UserIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.GetAll()); err != nil {
		panic(err)
	}
})

// UserShow handles GET /users/:userID
var UserShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := types.UUID(vars["userID"])
	user := user.FindByID(userID)
	if len(user.ID) > 0 {
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

// UserCreate handles POST /users/:userID
var UserCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"email"`
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	u := user.User{}
	u.Username = input.Username
	u.Password = input.Password
	u.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}
	if err := user.Validate(u); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := user.FindByUsername(u.Username); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := user.FindByAddress(u.Address); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Address is taken.")
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Email is invalid.")
		return
	}

	u = user.Create(u)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		fmt.Fprint(w, err)
		panic(err)
	}
})

// UserUpdate handles PUT /users/:userID
var UserUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"email"`
	}

	vars := mux.Vars(r)
	userID := types.UUID(vars["userID"])

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	u := user.User{}
	u.ID = userID
	u.Username = input.Username
	u.Password = input.Password
	u.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}

	if err := user.Validate(u); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := user.FindByID(u.ID); len(findUser.ID) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := user.FindByUsername(u.Username); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := user.FindByAddress(u.Address); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Address is taken.")
		return
	}

	u = user.Update(u)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
})

// UserPatch handles PATCH /users/:userID
var UserPatch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"email"`
	}

	vars := mux.Vars(r)
	userID := types.UUID(vars["userID"])

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	u := user.User{}
	u.ID = userID
	u.Username = input.Username
	u.Password = input.Password
	u.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}
	if err := user.Validate(u); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := user.FindByID(u.ID); len(findUser.ID) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := user.FindByUsername(u.Username); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := user.FindByAddress(u.Address); len(findUser.ID) > 0 && findUser.ID != u.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Address is taken.")
		return
	}

	u = user.Patch(u)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		panic(err)
	}
})

// UserDelete handles DELETE /users/:userID
var UserDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := types.UUID(vars["userID"])

	if err := user.Delete(userID); err == nil {
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
