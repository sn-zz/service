// Package user manages the users for the application.
//
// sn - https://github.com/sn
package user

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"github.com/sn/service/helpers"
	"github.com/sn/service/types"
)

// User represents a user
type User struct {
	ID       types.UUID
	Username string
	Password string
	Address  *mail.Address
	Created  time.Time
	Updated  time.Time
}

var users []User

// CheckPassword validates a password
func CheckPassword(u User, password string) bool {
	return u.Password == helpers.GeneratePasswordHash(password)
}

// GetAll returns all users
func GetAll() []User {
	return users
}

// FindByID looks for a user given a UUID
func FindByID(id types.UUID) User {
	for _, u := range users {
		if u.ID == id {
			return u
		}
	}
	return User{}
}

// FindByAddress finds a user by address
func FindByAddress(address *mail.Address) User {
	for _, u := range users {
		if u.Address.Address == address.Address {
			return u
		}
	}
	return User{}
}

// FindByUsername finds a user by username
func FindByUsername(username string) User {
	for _, u := range users {
		if u.Username == username {
			return u
		}
	}
	return User{}
}

// Validate validates a username, password, and email
//
// A user is valid if:
// - the username contains only alphanumerical characters,
// - the password
//   - is longer than 10 characters,
//   - contains at least one digit,
//   - contains at least one lowercase letter,
//   - contains at least one uppercase letter,
//   - contains at least one special character.
func Validate(user User) error {
	usernameRegex := regexp.MustCompile(`^[[:alnum:]]+$`)
	if len(user.Username) > 0 && !usernameRegex.MatchString(user.Username) {
		return fmt.Errorf("Username is invalid.")
	}

	if len(user.Password) > 0 {
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
	}
	return nil
}

// Create adds a user to the users list
func Create(user User) User {
	user.ID = helpers.GenerateUUID()
	user.Password = helpers.GeneratePasswordHash(user.Password)
	user.Created = time.Now()
	users = append(users, user)
	return user
}

// Update updates a user in the users list based on the user ID
func Update(user User) User {
	user.Password = helpers.GeneratePasswordHash(user.Password)
	for i, u := range users {
		if u.ID == user.ID {
			user.Updated = time.Now()
			users[i] = user
			return users[i]
		}
	}
	return User{}
}

// Patch patches a user in the users list based on the user ID
func Patch(user User) User {
	for i, u := range users {
		if u.ID == user.ID {
			if user.Address.Address != "" {
				u.Address = user.Address
			}
			if user.Username != "" {
				u.Username = user.Username
			}
			if user.Password != "" {
				u.Password = helpers.GeneratePasswordHash(user.Password)
			}
			u.Updated = time.Now()
			users[i] = u
			return users[i]
		}
	}
	return User{}
}

// Delete deletes a user based on the user ID
func Delete(id types.UUID) error {
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Not found")
}
