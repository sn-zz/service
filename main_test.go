// sn - https://github.com/sn
package main

import (
	"net/mail"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	address, err := mail.ParseAddress("zg@zk.gd")
	if err != nil {
		panic(err)
	}
	users = append(users, User{Id: GenerateUuid(), Username: "zg", Password: GeneratePasswordHash("s3cr3t"), Address: address, Created: time.Now()})
	address, err = mail.ParseAddress("bob@zk.gd")
	if err != nil {
		panic(err)
	}
	users = append(users, User{Id: GenerateUuid(), Username: "bob", Password: GeneratePasswordHash("s3cr3t"), Address: address, Created: time.Now()})

	sessions = append(sessions, Session{Id: GenerateUuid(), UserId: users[0].Id, Expires: time.Now().Add(SessionTime)}, Session{Id: GenerateUuid(), UserId: users[1].Id, Expires: time.Now().Add(SessionTime)})
	os.Exit(m.Run())
}
