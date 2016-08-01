// sn - https://github.com/sn
package session

import (
	"fmt"
	"time"

	"github.com/sn/service/helpers"
	"github.com/sn/service/types"
)

// Session contains a user's session
type Session struct {
	ID      types.UUID
	UserID  types.UUID
	Expires time.Time
}

// Expiration represents how much time a session lasts (one day)
const Expiration = 86400

var sessions []Session

// Create creates a new session
func Create(userID types.UUID) Session {
	s := Session{ID: helpers.GenerateUUID(), UserID: userID, Expires: time.Now().Add(Expiration)}
	sessions = append(sessions, s)
	return s
}

// Get retrieves a session given a user ID
func Get(id types.UUID) Session {
	for _, s := range sessions {
		if s.ID == id {
			return s
		}
	}
	return Session{}
}

// GetAll retrieves all sessions
func GetAll() []Session {
	return sessions
}

// Expire sets the expiration of a session well into the past
func Expire(id types.UUID) error {
	for i, s := range sessions {
		if s.ID == id {
			sessions[i].Expires = time.Time{}
			return nil
		}
	}
	return fmt.Errorf("Could not find session")
}

// Find retrieves a session given a session hash
func Find(hash string) Session {
	for _, s := range sessions {
		if helpers.GenerateSha1Hash(string(s.ID)) == hash {
			return s
		}
	}
	return Session{}
}

// Bump bumps the expiration time up for a given session UUID
func Bump(id types.UUID) error {
	for i, s := range sessions {
		if s.ID == id {
			sessions[i].Expires = time.Now().Add(Expiration)
			return nil
		}
	}
	return fmt.Errorf("Could not find session")
}

// Clean removes any expired sessions
func Clean() {
	for i := 0; i < len(sessions); i++ {
		if time.Now().After(sessions[i].Expires) {
			Remove(sessions[i].ID)
			i--
		}
	}
}

// Remove removes a session from the existing sessions
func Remove(id types.UUID) error {
	for i, s := range sessions {
		if s.ID == id {
			sessions = append(sessions[:i], sessions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find session")
}
