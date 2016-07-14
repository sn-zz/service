// sn - https://github.com/sn
package main

import "fmt"

type User struct {
    Id        string    `json:"userId"`
    Username  string    `json:"username"`
    Password  string    `json:"password"`
    Email     string    `json:"email"`
}

type Users []User

var users Users

// Give us some seed data
func init() {
    CreateUser(User{Username:"zg",Password:"s3cr3t",Email:"zg@zk.gd"})
    CreateUser(User{Username:"bob",Password:"s3cr3t",Email:"bob@zk.gd"})
}

func IsEmailTaken(email string) bool {
    user := FindUserByEmail(email)
    return len(user.Id) > 0
}

func IsUsernameTaken(username string) bool {
    user := FindUserByUsername(username)
    return len(user.Id) > 0
}

func CheckPassword(u User, password string) bool {
    return u.Password == GenerateHash(password)
}

func FindUserById(id string) User {
    for _, u := range users {
        if u.Id == id {
            return u
        }
    }
    return User{}
}

func FindUserByEmail(email string) User {
    for _, u := range users {
        if u.Email == email {
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

//this is bad, I don't think it passes race conditions
func CreateUser(u User) (User, error) {
    if u.Id = GenerateUuid(); u.Id != "" {
        u.Password = GenerateHash(u.Password)
        users = append(users, u)
        return u, nil
    }
    return User{}, fmt.Errorf("Could not generate UUID")
}

func UpdateUser(id string, user User) (User, error) {
    for _, u := range users {
        if u.Id == id {
            u = user
            return u, nil
        }
    }
    return User{}, fmt.Errorf("Not found")
}

func PatchUser(id string, u User) (User, error) {
    for i, user := range users {
        if user.Id == id {
            if u.Email != "" {
                user.Email = u.Email
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

func DeleteUser(id string) error {
    for i, u := range users {
        if u.Id == id {
            users = append(users[:i], users[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Not found")
}
