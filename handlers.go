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
		if session := FindSession(auth[0]); session.Id != "" {
			u := FindUserById(session.UserId)
			fmt.Fprintf(w, "Welcome, %s!\n", u.Username)
			err := UpdateSessionTime(session.Id)
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
	refUser := FindUserById(user.Id)
	if len(refUser.Id) > 0 {
		if CheckPassword(refUser, user.Password) {
			w.WriteHeader(http.StatusOK)
			s := CreateSession(refUser.Id)
			fmt.Fprintf(w, "%s", GenerateSha1Hash(string(s.Id)))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("userId %s, password %s", user.Id, user.Password)

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

// UserShow handles GET /users/:userId
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

// UserCreate handles POST /users/:userId
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
	if findUser := FindUserByUsername(user.Username); len(findUser.Id) > 0 && findUser.Id != user.Id {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.Id) > 0 && findUser.Id != user.Id {
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

// UserUpdate handles PUT /users/:userId
var UserUpdate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"email"`
	}

	vars := mux.Vars(r)
	userId := uuid(vars["userId"])

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
	user.Id = userId
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
	if findUser := FindUserById(user.Id); len(findUser.Id) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := FindUserByUsername(user.Username); len(findUser.Id) > 0 && findUser.Id != user.Id {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.Id) > 0 && findUser.Id != user.Id {
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

// UserPatch handles PATCH /users/:userId
var UserPatch = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"email"`
	}

	vars := mux.Vars(r)
	userId := uuid(vars["userId"])

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
	user.Id = userId
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
	if findUser := FindUserById(user.Id); len(findUser.Id) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not found")
		return
	}
	if findUser := FindUserByUsername(user.Username); len(findUser.Id) > 0 && findUser.Id != user.Id {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Username is taken.")
		return
	}
	if findUser := FindUserByAddress(user.Address); len(findUser.Id) > 0 && findUser.Id != user.Id {
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

// UserDelete handles DELETE /users/:userId
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
