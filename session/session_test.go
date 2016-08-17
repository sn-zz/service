// Package session manages the sessions for the application.
//
// sn - https://github.com/sn
package session

import (
	"log"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sn/service/helpers"
	"github.com/sn/service/user"
)

func TestCreate(t *testing.T) {
	userID := helpers.GenerateUUID()
	if userID == "" {
		t.Error("Could not generate UUID")
	}
	newSession := Create(userID)
	if newSession.UserID != userID {
		t.Error("User UUID mismatch.")
	}
	if newSession.Expires.IsZero() {
		t.Error("Expiration was not properly set.")
	}
}

func TestGet(t *testing.T) {
	s := Session{}
	s = Get(s.ID)
	if len(s.ID) != 0 {
		t.Error("Get fail should return empty session.")
	}
	sessions := GetAll()
	s = Get(sessions[0].ID)
	if s.Expires.Sub(sessions[0].Expires) != 0 {
		t.Error("Incorrect session was obtained.")
	}
	users := user.GetAll()
	if s.UserID != users[0].ID {
		t.Error("Incorrect user ID associated with session.")
	}
}

func TestGetAll(t *testing.T) {
	sessions := GetAll()
	if len(sessions) == 0 {
		t.Error("Incorrect sessions length.")
	}
}

func TestExpire(t *testing.T) {
	s := Session{}
	err := Expire(s.ID)
	if err.Error() != "Could not find session" {
		t.Error("Remove fail should specify not found.")
	}
	sessions := GetAll()
	for _, s := range sessions {
		Expire(s.ID)
	}
	for _, s := range sessions {
		if !s.Expires.IsZero() {
			t.Error("Incorrect session expiration.")
		}
	}
}

func TestFind(t *testing.T) {
	s := Session{}
	s = Find("")
	if len(s.ID) != 0 {
		t.Error("Find fail should return empty session.")
	}
	sessions := GetAll()
	sessionHash := helpers.GenerateSha1Hash(string(sessions[0].ID))
	s = Find(sessionHash)
	if s.Expires.Sub(sessions[0].Expires) != 0 {
		t.Error("Incorrect session was obtained.")
	}
	users := user.GetAll()
	if s.UserID != users[0].ID {
		t.Error("Incorrect user ID associated with session.")
	}
}

func TestBump(t *testing.T) {
	s := Session{}
	err := Bump(s.ID)
	if err.Error() != "Could not find session" {
		t.Error("Bump fail should specify not found.")
	}
	sessions := GetAll()
	s = Get(sessions[0].ID)
	err = Bump(sessions[0].ID)
	if err != nil {
		t.Error(err)
	}
	if s.Expires.After(Get(sessions[0].ID).Expires) {
		t.Error("Expiration was not updated.")
	}
}

func TestClean(t *testing.T) {
	sessions := GetAll()
	for _, s := range sessions {
		Expire(s.ID)
	}
	Clean()
	sessions = GetAll()
	if len(sessions) != 0 {
		t.Error("Sessions did not clean correctly.")
	}
}

func TestRemove(t *testing.T) {
	s := Session{}
	err := Remove(s.ID)
	if err.Error() != "Could not find session" {
		t.Error("Remove fail should specify not found.")
	}
	s = Create(helpers.GenerateUUID())
	sessions := GetAll()
	err = Remove(s.ID)
	if len(sessions) == len(GetAll()) {
		t.Error("Session was not removed.")
	}
}

func TestMain(m *testing.M) {
	usernames := [4]string{"alex", "blake", "corey", "devon"}
	for _, un := range usernames {
		addr, err := mail.ParseAddress(strings.Title(un) + "<" + un + "@example.com>")
		if err != nil {
			log.Fatal(err)
		}
		u := user.User{Username: un, Password: helpers.GeneratePasswordHash("s3cr3t"), Address: addr, Created: time.Now()}
		u = user.Create(u)
		Create(u.ID)
	}

	os.Exit(m.Run())
}
