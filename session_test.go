// sn - https://github.com/sn
package main

import (
	"testing"
	"time"
)

func TestCreateSession(t *testing.T) {
	userId := GenerateUuid()
	if userId == "" {
		t.Errorf("Could not generate UUID")
	}
	newSession := CreateSession(userId)
	if newSession.UserId != userId {
		t.Errorf("User UUID mismatch.")
	}
	if newSession.Expires.IsZero() {
		t.Errorf("Expiration was not properly set.")
	}
}

func TestGetSession(t *testing.T) {
	session := GetSession(sessions[0].Id)
	if session.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	if session.UserId != users[0].Id {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestFindSession(t *testing.T) {
	sessionHash := GenerateSha1Hash(string(sessions[0].Id))
	session := FindSession(sessionHash)
	if session.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	if session.UserId != users[0].Id {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestUpdateSessionTime(t *testing.T) {
	session := GetSession(sessions[0].Id)
	err := UpdateSessionTime(sessions[0].Id)
	if err != nil {
		panic(err)
	}
	if session.Expires.After(GetSession(sessions[0].Id).Expires) {
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
