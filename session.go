// sn - https://github.com/sn
package main

import (
	"fmt"
	"time"
)

// Session contains a user's session
type Session struct {
	Id      uuid      `json:"sessionId"`
	UserId  uuid      `json:"userId"`
	Expires time.Time `json:"expires"`
}

// Sessions contains all sessions
type Sessions []Session

// SessionTime represents how much time a session lasts (one day)
const SessionTime = 86400

var sessions Sessions

// CreateSession creates a new session
func CreateSession(userId uuid) (Session, error) {
	if id := GenerateUuid(); id != "" {
		s := Session{Id: id, UserId: userId, Expires: time.Now().Add(SessionTime)}
		sessions = append(sessions, s)
		return s, nil
	}
	return Session{}, fmt.Errorf("Could not generate UUID")
}

// GetSession retrieves a session given a user ID
func GetSession(id uuid) Session {
	for _, s := range sessions {
		if s.Id == id {
			return s
		}
	}
	return Session{}
}

// FindSession retrieves a session given a session hash
func FindSession(hash string) Session {
	for _, s := range sessions {
		if GenerateSha1Hash(string(s.Id)) == hash {
			return s
		}
	}
	return Session{}
}

// UpdateSessionTime updates a given session UUID
func UpdateSessionTime(id uuid) error {
	for i, s := range sessions {
		if s.Id == id {
			sessions[i].Expires = time.Now().Add(SessionTime)
		}
	}
	return fmt.Errorf("Could not find session")
}

// CleanSessions removes any expired sessions
func CleanSessions() {
	for i, s := range sessions {
		if time.Now().After(s.Expires) {
			sessions = append(sessions[:i], sessions[i+1:]...)
		}
	}
}
