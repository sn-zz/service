// sn - https://github.com/sn
package main

import (
    "fmt"
    "net/mail"
    "regexp"
    "time"
)

// User represents a user
type User struct {
    Id       uuid
    Username string
    Password string
    Address  *mail.Address
    Created  time.Time
    Updated  time.Time
}

// Users contains all users
type Users []User

var users Users

// IsAddressTaken checks if an address is taken
func IsAddressTaken(address string) bool {
    Address, err := mail.ParseAddress(address)
    if err != nil {
        panic(err)
    }
    user := FindUserByAddress(Address)
    return len(user.Id) > 0
}

// IsUsernameTaken checks if a username is taken
func IsUsernameTaken(username string) bool {
    user := FindUserByUsername(username)
    return len(user.Id) > 0
}

// CheckPassword validates a password
func CheckPassword(u User, password string) bool {
    return u.Password == GeneratePasswordHash(password)
}

// FindUserById looks for a user given a UUID
func FindUserById(id uuid) User {
    for _, u := range users {
        if u.Id == id {
            return u
        }
    }
    return User{}
}

// FindUserByAddress finds a user by address
func FindUserByAddress(address *mail.Address) User {
    for _, u := range users {
        if u.Address.Address == address.Address {
            return u
        }
    }
    return User{}
}

// FindUserByUsername finds a user by username
func FindUserByUsername(username string) User {
    for _, u := range users {
        if u.Username == username {
            return u
        }
    }
    return User{}
}

// ValidateUser validates a username, password, and email
//
// A user is valid if:
// - the username contains only alphanumerical characters,
// - the password
//   - is longer than 10 characters,
//   - contains at least one digit,
//   - contains at least one lowercase letter,
//   - contains at least one uppercase letter,
//   - contains at least one special character.
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

// CreateUser adds a user to the users list
func CreateUser(user User) (User, error) {
    if user.Id = GenerateUuid(); user.Id != "" {
        user.Password = GeneratePasswordHash(user.Password)
        user.Created = time.Now()
        users = append(users, user)
        return user, nil
    }
    return User{}, fmt.Errorf("Could not generate UUID")
}

// UpdateUser updates a user in the users list based on the user ID
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

// PatchUser patches a user in the users list based on the user ID
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

// DeleteUser deletes a user based on the user ID
func DeleteUser(id uuid) error {
    for i, u := range users {
        if u.Id == id {
            users = append(users[:i], users[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Not found")
}
