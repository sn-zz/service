// sn - https://github.com/sn
package main

import (
    "fmt"
    "net/mail"
    "time"
)

type User struct {
    Id       uuid          `json:"userId"`
    Username string        `json:"username"`
    Password string        `json:"password"`
    Address  *mail.Address `json:"email"`
    Created  time.Time     `json:"created"`
    Updated  time.Time     `json:"updated"`
}

type Users []User

var users Users

// Give us some seed data
func init() {
    address, _ := mail.ParseAddress("zg@zk.gd")
    CreateUser(User{Username:"zg",Password:"s3cr3t",Address:address,Created:time.Now()})
    address, _ = mail.ParseAddress("zg@zk.gd")
    CreateUser(User{Username:"bob",Password:"s3cr3t",Address:address,Created:time.Now()})
}

func IsAddressTaken(address string) bool {
    Address, _ := mail.ParseAddress(address)
    user := FindUserByAddress(Address)
    return len(user.Id) > 0
}

func IsUsernameTaken(username string) bool {
    user := FindUserByUsername(username)
    return len(user.Id) > 0
}

func CheckPassword(u User, password string) bool {
    return u.Password == GeneratePasswordHash(password)
}

func FindUserById(id uuid) User {
    for _, u := range users {
        if u.Id == id {
            return u
        }
    }
    return User{}
}

func FindUserByAddress(address *mail.Address) User {
    for _, u := range users {
        if u.Address == address {
            return u
        }
    }
    return User{}
}

func FindUserByUsername(username string) User {
    for _, u := range users {
        if u.Username == username {
            return u
        }
    }
    return User{}
}

func CreateUser(u User) (User, error) {
    if u.Id = GenerateUuid(); u.Id != "" {
        u.Password = GeneratePasswordHash(u.Password)
        users = append(users, u)
        return u, nil
    }
    return User{}, fmt.Errorf("Could not generate UUID")
}

func UpdateUser(id uuid, user User) (User, error) {
    for _, u := range users {
        if u.Id == id {
            u = user
            return u, nil
        }
    }
    return User{}, fmt.Errorf("Not found")
}

func PatchUser(id uuid, u User) (User, error) {
    for i, user := range users {
        if user.Id == id {
            if u.Address.Address != "" {
                user.Address = u.Address
            }
            if u.Username != "" {
                user.Username = u.Username
            }
            if u.Password != "" {
                user.Password = u.Password
            }
            users[i] = user
            return user, nil
        }
    }
    return User{}, fmt.Errorf("Not found")
}

func DeleteUser(id uuid) error {
    for i, u := range users {
        if u.Id == id {
            users = append(users[:i], users[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Not found")
}
