// sn - https://github.com/sn
package main

import (
    "fmt"
    "net/mail"
    "regexp"
    "time"
)

type User struct {
    Id       uuid
    Username string
    Password string
    Address  *mail.Address
    Created  time.Time
    Updated  time.Time
}

type Users []User

var users Users

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
        if u.Address.Address == address.Address {
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

func ValidateUser(user User) error {
    usernameRegex := regexp.MustCompile(`^[[:alnum:]]+$`)
    if !usernameRegex.MatchString(user.Username) {
        return fmt.Errorf("Username is invalid.")
    }

    length := regexp.MustCompile(`.{10,}`)
    digits := regexp.MustCompile(`[[:digit:]]`)
    lowers := regexp.MustCompile(`[[:lower:]]`)
    uppers := regexp.MustCompile(`[[:upper:]]`)
    special := regexp.MustCompile(`[!"#$%&'()*+,\-./:;<=>?@[\\\]^_{|}~\x60]`) // \x60 == `

    if !length.MatchString(user.Password) {
        return fmt.Errorf("Password must be 10 characters or longer.")
    }
    if !digits.MatchString(user.Password) {
        return fmt.Errorf("Password must contain a number.")
    }
    if !lowers.MatchString(user.Password) {
        return fmt.Errorf("Password must contain a lowercase letter.")
    }
    if !uppers.MatchString(user.Password) {
        return fmt.Errorf("Password must contain an uppercase letter.")
    }
    if !special.MatchString(user.Password) {
        return fmt.Errorf("Password must contain a special character.")
    }
    return nil
}

func CreateUser(user User) (User, error) {
    if user.Id = GenerateUuid(); user.Id != "" {
        user.Password = GeneratePasswordHash(user.Password)
        user.Created = time.Now()
        users = append(users, user)
        return user, nil
    }
    return User{}, fmt.Errorf("Could not generate UUID")
}

func UpdateUser(user User) (User, error) {
    user.Password = GeneratePasswordHash(user.Password)
    for i, u := range users {
        if u.Id == user.Id {
            user.Updated = time.Now()
            users[i] = user
            return users[i], nil
        }
    }
    return User{}, fmt.Errorf("Not found")
}

func PatchUser(user User) (User, error) {
    for i, u := range users {
        if u.Id == user.Id {
            if user.Address.Address != "" {
                u.Address = user.Address
            }
            if user.Username != "" {
                u.Username = user.Username
            }
            if user.Password != "" {
                u.Password = GeneratePasswordHash(user.Password)
            }
            u.Updated = time.Now()
            users[i] = u
            return users[i], nil
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
