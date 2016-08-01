// sn - https://github.com/sn
package main

import (
	"fmt"
	"time"
)

// Session contains a user's session
type Session struct {
	ID      uuid
	UserID  uuid
	Expires time.Time
}

// Sessions contains all sessions
type Sessions []Session

// SessionTime represents how much time a session lasts (one day)
const SessionTime = 86400

var sessions Sessions

// CreateSession creates a new session
func CreateSession(userID uuid) Session {
	s := Session{ID: GenerateUUID(), UserID: userID, Expires: time.Now().Add(SessionTime)}
	sessions = append(sessions, s)
	return s
}

// GetSession retrieves a session given a user ID
func GetSession(id uuid) Session {
	for _, s := range sessions {
		if s.ID == id {
			return s
		}
	}
	return Session{}
}

// FindSession retrieves a session given a session hash
func FindSession(hash string) Session {
	for _, s := range sessions {
		if GenerateSha1Hash(string(s.ID)) == hash {
			return s
		}
	}
	return Session{}
}

// UpdateSessionTime updates a given session UUID
func UpdateSessionTime(id uuid) error {
	for i, s := range sessions {
		if s.ID == id {
			sessions[i].Expires = time.Now().Add(SessionTime)
			return nil
		}
	}
	return fmt.Errorf("Could not find session")
}

// CleanSessions removes any expired sessions
func CleanSessions() {
	for i, s := range sessions {
		if time.Now().After(s.Expires) {
			if len(sessions) > 1 {
				sessions = append(sessions[:i], sessions[i+1:]...)
			} else {
				sessions = make([]Session, 0)
			}
		}
	}
}
