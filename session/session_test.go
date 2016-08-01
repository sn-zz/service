// sn - https://github.com/sn
package session

import (
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
		t.Errorf("Could not generate UUID")
	}
	newSession := Create(userID)
	if newSession.UserID != userID {
		t.Errorf("User UUID mismatch.")
	}
	if newSession.Expires.IsZero() {
		t.Errorf("Expiration was not properly set.")
	}
}

func TestGet(t *testing.T) {
	sessions := GetAll()
	s := Get(sessions[0].ID)
	if s.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	users := user.GetAll()
	if s.UserID != users[0].ID {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestGetAll(t *testing.T) {
	sessions := GetAll()
	if len(sessions) == 0 {
		t.Errorf("Incorrect sessions length.")
	}
}

func TestExpire(t *testing.T) {
	sessions := GetAll()
	for _, s := range sessions {
		Expire(s.ID)
	}
	for _, s := range sessions {
		if !s.Expires.IsZero() {
			t.Errorf("Incorrect session expiration.")
		}
	}
}

func TestFind(t *testing.T) {
	sessions := GetAll()
	sessionHash := helpers.GenerateSha1Hash(string(sessions[0].ID))
	s := Find(sessionHash)
	if s.Expires.Sub(sessions[0].Expires) != 0 {
		t.Errorf("Incorrect session was obtained.")
	}
	users := user.GetAll()
	if s.UserID != users[0].ID {
		t.Errorf("Incorrect user ID associated with session.")
	}
}

func TestBump(t *testing.T) {
	sessions := GetAll()
	s := Get(sessions[0].ID)
	err := Bump(sessions[0].ID)
	if err != nil {
		panic(err)
	}
	if s.Expires.After(Get(sessions[0].ID).Expires) {
		t.Errorf("Expiration was not updated.")
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
		t.Errorf("Sessions did not clean correctly.")
	}
}

func TestRemove(t *testing.T) {
	s := Create(helpers.GenerateUUID())
	sessions := GetAll()
	Remove(s.ID)
	if len(sessions) == len(GetAll()) {
		t.Errorf("Session was not removed.")
	}
}

func TestMain(m *testing.M) {
	usernames := [4]string{"alex", "blake", "corey", "devon"}
	for _, un := range usernames {
		addr, err := mail.ParseAddress(strings.Title(un) + "<" + un + "@example.com>")
		if err != nil {
			panic(err)
		}
		u := user.User{Username: un, Password: helpers.GeneratePasswordHash("s3cr3t"), Address: addr, Created: time.Now()}
		u = user.Create(u)
		Create(u.ID)
	}

	os.Exit(m.Run())
}
