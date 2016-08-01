// sn - https://github.com/sn
package main

import (
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	userID := GenerateUUID()
	if userID == "" {
		t.Errorf("Could not generate UUID")
	}
	newSession := CreateSession(userID)
	if newSession.UserID != userID {
		t.Errorf("User UUID mismatch.")
	}
	if newSession.Expires.IsZero() {
		t.Errorf("Expiration was not properly set.")
	}
}

func TestGetSession(t *testing.T) {
	session := GetSession(sessions[0].ID)
	if session.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	if session.UserID != users[0].ID {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestFindSession(t *testing.T) {
	sessionHash := GenerateSha1Hash(string(sessions[0].ID))
	session := FindSession(sessionHash)
	if session.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	if session.UserID != users[0].ID {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestUpdateSessionTime(t *testing.T) {
	session := GetSession(sessions[0].ID)
	err := UpdateSessionTime(sessions[0].ID)
	if err != nil {
		panic(err)
	}
	if session.Expires.After(GetSession(sessions[0].ID).Expires) {
		t.Errorf("Expiration was not updated.")
	}
}

func TestCleanSessions(t *testing.T) {
	for i := range sessions {
		sessions[i].Expires = time.Now()
	}
	CleanSessions()
	if len(sessions) > 0 {
		t.Errorf("Sessions did not clean correctly.")
	}
}
