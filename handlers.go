// sn - https://github.com/sn
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
)

// Index handles GET /index
var Index = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if auth := r.Header["Authorization"]; auth != nil {
		if session := FindSession(auth[0]); session.ID != "" {
			u := FindUserByID(session.UserID)
			fmt.Fprintf(w, "Welcome, %s!\n", u.Username)
			err := UpdateSessionTime(session.ID)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	fmt.Fprint(w, "Welcome!\n")
})

// Auth handles POST /auth
var Auth = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := User{}
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
	refUser := FindUserByID(user.ID)
	if len(refUser.ID) > 0 {
		if CheckPassword(refUser, user.Password) {
			w.WriteHeader(http.StatusOK)
			s := CreateSession(refUser.ID)
			fmt.Fprintf(w, "%s", GenerateSha1Hash(string(s.ID)))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("userID %s, password %s", user.ID, user.Password)

	w.WriteHeader(http.StatusNotFound)
})

// UserIndex handles GET /users
var UserIndex = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
})

// UserShow handles GET /users/:userID
var UserShow = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := uuid(vars["userID"])
	user := FindUserByID(userID)
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

	user := User{}
	user.Username = input.Username
	user.Password = input.Password
	user.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}
	if err := ValidateUser(user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := FindUserByUsername(user.Username); len(findUser.ID) > 0 && findUser.ID != user.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.ID) > 0 && findUser.ID != user.ID {
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

	user = CreateUser(user)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
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
	userID := uuid(vars["userID"])

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

	user := User{}
	user.ID = userID
	user.Username = input.Username
	user.Password = input.Password
	user.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}

	if err := ValidateUser(user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := FindUserByID(user.ID); len(findUser.ID) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := FindUserByUsername(user.Username); len(findUser.ID) > 0 && findUser.ID != user.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.ID) > 0 && findUser.ID != user.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Address is taken.")
		return
	}

	user = UpdateUser(user)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(user); err != nil {
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
	userID := uuid(vars["userID"])

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

	user := User{}
	user.ID = userID
	user.Username = input.Username
	user.Password = input.Password
	user.Address, err = mail.ParseAddress(input.Address)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to parse address.")
		return
	}
	if err := ValidateUser(user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	if findUser := FindUserByID(user.ID); len(findUser.ID) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := FindUserByUsername(user.Username); len(findUser.ID) > 0 && findUser.ID != user.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.ID) > 0 && findUser.ID != user.ID {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Address is taken.")
		return
	}

    user = PatchUser(user)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(user); err != nil {
        panic(err)
    }
})

// UserDelete handles DELETE /users/:userID
var UserDelete = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := uuid(vars["userID"])

	if err := DeleteUser(userID); err == nil {
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
