// Package user manages the users for the application.
//
// sn - https://github.com/sn
package user

import (
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sn/service/helpers"
)

func TestCheckPassword(t *testing.T) {
	users := GetAll()
	u := FindByID(users[0].ID)
	correctPassword := "1@E4s67890"
	incorrectPassword := "s3cret"

	if !CheckPassword(u, correctPassword) {
		t.Error("Expected password success, got failure.")
	}

	if CheckPassword(u, incorrectPassword) {
		t.Error("Expected password failure, got success.")
	}
}

func TestGetAll(t *testing.T) {
	users := GetAll()
	if len(users) == 0 {
		t.Error("Incorrect users length.")
	}
}

func TestFindByID(t *testing.T) {
	users := GetAll()
	knownID := users[0].ID
	unknownID := helpers.GenerateUUID()

	if u := FindByID(knownID); len(u.ID) == 0 {
		t.Error("Expected known user ID, got unknown user ID.")
	}

	if u := FindByID(unknownID); len(u.ID) > 0 {
		t.Error("Expected unknown user ID, got known user ID.")
	}
}

func TestFindByAddress(t *testing.T) {
	users := GetAll()
	knownAddress, err := mail.ParseAddress(users[0].Address.Address)
	if err != nil {
		t.Error(err)
	}
	unknownAddress, err := mail.ParseAddress("test@example.com")
	if err != nil {
		t.Error(err)
	}

	if u := FindByAddress(knownAddress); len(u.ID) == 0 {
		t.Error("Expected known address, got unknown address.")
	}

	if u := FindByAddress(unknownAddress); len(u.ID) > 0 {
		t.Error("Expected unknown address, got known address.")
	}
}

func TestFindByUsername(t *testing.T) {
	users := GetAll()
	knownUsername := users[0].Username
	unknownUsername := "unknown-username"

	if u := FindByUsername(knownUsername); len(u.ID) == 0 {
		t.Error("Expected known user, got unknown user.")
	}

	if u := FindByUsername(unknownUsername); len(u.ID) > 0 {
		t.Error("Expected unknown user, got known user.")
	}
}

func TestValidate(t *testing.T) {
	address, err := mail.ParseAddress("test@example.com")
	if err != nil {
		t.Error(err)
	}
	u := User{Username: "", Password: "123456789", Address: address}
	if err := Validate(u); err == nil { // Complains about username length == 0
		t.Error(err)
	}
	u.Username = "@@"
	if err := Validate(u); err == nil { // Complains about illegal characters
		t.Error(err)
	}
	u.Username = "zg"
	if err := Validate(u); err == nil { // Complains about length
		t.Error(err)
	}
	u.Password = "aaaaaaaaaa"
	if err := Validate(u); err == nil { // Complains about digits
	}
	u.Password = "0123456789"
	if err := Validate(u); err == nil { // Complains about lowercase
		t.Error(err)
	}
	u.Password = "01234s6789"
	if err := Validate(u); err == nil { // Complains about uppercase
		t.Error(err)
	}
	u.Password = "01z34S6789"
	if err := Validate(u); err == nil { // Complains about special characters
		t.Error(err)
	}
	u.Password = "@1z34S6789"
	if err := Validate(u); err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	users := GetAll()
	password := "S3crET!@#$"
	address, err := mail.ParseAddress(users[0].Address.Address + ".com")
	if err != nil {
		t.Error(err)
	}
	currentUserCount := len(users)
	u := Create(User{Username: "zzg", Password: password, Address: address})
	if len(GetAll()) == currentUserCount {
		t.Error("User wasn't created.")
	}
	if u.Created.IsZero() {
		t.Error("User creation time not set.")
	}
	if !CheckPassword(u, password) {
		t.Error("User password incorrectly set.")
	}
}

func TestUpdate(t *testing.T) {
	users := GetAll()
	address, err := mail.ParseAddress("zg@zk.gd")
	if err != nil {
		t.Error(err)
	}
	u := Update(User{})
	if len(u.ID) != 0 {
		t.Error("Patch fail should return empty user.")
	}
	userBeforeUpdate := FindByID(users[0].ID)
	updatedUser := User{ID: users[0].ID, Username: "zgg", Password: "S3crET!@#$", Address: address}
	u = Update(updatedUser)
	if len(u.ID) == 0 {
		t.Error("User was not found.")
	}
	if u.Username != updatedUser.Username {
		t.Error("Username was not updated.")
	}
	if !CheckPassword(u, updatedUser.Password) {
		t.Error("Password was not updated.")
	}
	if u.Address.Address != updatedUser.Address.Address {
		t.Error("Email was not updated.")
	}
	if u.Updated.Sub(userBeforeUpdate.Updated) == 0 {
		t.Error("Last Updated not updated.")
	}
}

func TestPatch(t *testing.T) {
	users := GetAll()
	address, err := mail.ParseAddress("zzg@zk.gd")
	if err != nil {
		t.Error(err)
	}
	u := Patch(User{})
	if len(u.ID) != 0 {
		t.Error("Patch fail should return empty user.")
	}
	userToPatch := FindByID(users[0].ID)
	userToPatch.Username = "zzg"
	userToPatch.Password = "S3crET!@#$"
	userToPatch.Address = address
	u = Patch(userToPatch)
	if len(u.ID) == 0 {
		t.Error("User was not found.")
	}
	if u.Username != userToPatch.Username {
		t.Error("Username was not patched.")
	}
	if !CheckPassword(u, userToPatch.Password) {
		t.Error("Password was not patched.")
	}
	if u.Address.Address != userToPatch.Address.Address {
		t.Error("Email was not patched.")
	}
	if u.Updated.Sub(userToPatch.Updated) == 0 {
		t.Error("Last Updated not patched.")
	}
}

func TestDelete(t *testing.T) {
	users := GetAll()
	u := User{}
	err := Delete(u.ID)
	if err.Error() != "Not found" {
		t.Error("Delete fail should return empty user.")
	}
	u = FindByID(users[0].ID)
	err = Delete(u.ID)
	if err != nil {
		t.Error(err)
	}
	if users[0].ID == u.ID {
		t.Error("User was not deleted.")
	}
}

func TestMain(m *testing.M) {
	usernames := [4]string{"alex", "blake", "corey", "devon"}
	for _, un := range usernames {
		addr, err := mail.ParseAddress(strings.Title(un) + "<" + un + "@example.com>")
		if err != nil {
			panic(err)
		}
		u := User{Username: un, Password: "1@E4s67890", Address: addr, Created: time.Now()}
		_ = Create(u)
	}

	os.Exit(m.Run())
}
