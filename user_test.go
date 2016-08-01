// sn - https://github.com/sn
package main

import (
	"net/mail"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	id := users[0].ID
	user := FindUserByID(id)
	correctPassword := "s3cr3t"
	incorrectPassword := "s3cret"

	if !CheckPassword(user, correctPassword) {
		t.Errorf("Expected password success, got failure.")
	}

	if CheckPassword(user, incorrectPassword) {
		t.Errorf("Expected password failure, got success.")
	}
}

func TestFindUserByID(t *testing.T) {
	knownID := users[0].ID
	unknownID := GenerateUUID()

	if user := FindUserByID(knownID); len(user.ID) == 0 {
		t.Errorf("Expected known user ID, got unknown user ID.")
	}

	if user := FindUserByID(unknownID); len(user.ID) > 0 {
		t.Errorf("Expected unknown user ID, got known user ID.")
	}
}

func TestFindUserByAddress(t *testing.T) {
	knownAddress, err := mail.ParseAddress(users[0].Address.Address)
	if err != nil {
		t.Error(err)
	}
	unknownAddress, err := mail.ParseAddress("test@example.com")
	if err != nil {
		t.Error(err)
	}

	if user := FindUserByAddress(knownAddress); len(user.ID) == 0 {
		t.Errorf("Expected known address, got unknown address.")
	}

	if user := FindUserByAddress(unknownAddress); len(user.ID) > 0 {
		t.Errorf("Expected unknown address, got known address.")
	}
}

func TestFindUserByUsername(t *testing.T) {
	knownUsername := users[0].Username
	unknownUsername := "unknown-username"

	if user := FindUserByUsername(knownUsername); len(user.ID) == 0 {
		t.Errorf("Expected known user, got unknown user.")
	}

	if user := FindUserByUsername(unknownUsername); len(user.ID) > 0 {
		t.Errorf("Expected unknown user, got known user.")
	}
}

func TestValidateUser(t *testing.T) {
	address, err := mail.ParseAddress("test@example.com")
	if err != nil {
		t.Error(err)
	}
	user := User{Username: "zg", Password: "123456789", Address: address}
	if err := ValidateUser(user); err == nil { // Complains about length
		t.Error(err)
	}
	user.Password = "0123456789"
	if err := ValidateUser(user); err == nil { // Complains about lowercase
		t.Error(err)
	}
	user.Password = "01234s6789"
	if err := ValidateUser(user); err == nil { // Complains about uppercase
		t.Error(err)
	}
	user.Password = "01z34S6789"
	if err := ValidateUser(user); err == nil { // Complains about special characters
		t.Error(err)
	}
	user.Password = "@1z34S6789"
	if err := ValidateUser(user); err != nil {
		t.Error(err)
	}
}

func TestCreateUser(t *testing.T) {
	password := "S3crET!@#$"
	address, err := mail.ParseAddress(users[0].Address.Address)
	if err != nil {
		t.Error(err)
	}
	newUser := User{Username: "zzg", Password: password, Address: address}
	u := CreateUser(newUser)
	if len(users) != 3 {
		t.Errorf("User wasn't created.")
	}
	if u.Created.IsZero() {
		t.Errorf("User creation time not set.")
	}
	if !CheckPassword(u, password) {
		t.Errorf("User password incorrectly set.")
	}
}

func TestUpdateUser(t *testing.T) {
	address, err := mail.ParseAddress("zg@zk.gd")
	if err != nil {
		t.Error(err)
	}
	userBeforeUpdate := FindUserByID(users[0].ID)
	updatedUser := User{ID: users[0].ID, Username: "zgg", Password: "S3crET!@#$", Address: address}
	user := UpdateUser(updatedUser)
    if len(user.ID) == 0 {
        t.Errorf("User was not found.")
    }
	if user.Username != updatedUser.Username {
		t.Errorf("Username was not updated.")
	}
	if !CheckPassword(user, updatedUser.Password) {
		t.Errorf("Password was not updated.")
	}
	if user.Address.Address != updatedUser.Address.Address {
		t.Errorf("Email was not updated.")
	}
	if user.Updated.Sub(userBeforeUpdate.Updated) == 0 {
		t.Errorf("Last Updated not updated.")
	}
}

func TestPatchUser(t *testing.T) {
	address, err := mail.ParseAddress("zzg@zk.gd")
	if err != nil {
		t.Error(err)
	}
	userToPatch := FindUserByID(users[0].ID)
	userToPatch.Username = "zzg"
	userToPatch.Password = "S3crET!@#$"
	userToPatch.Address = address
	user := PatchUser(userToPatch)
    if len(user.ID) == 0 {
        t.Errorf("User was not found.")
    }
	if user.Username != userToPatch.Username {
		t.Errorf("Username was not patched.")
	}
	if !CheckPassword(user, userToPatch.Password) {
		t.Errorf("Password was not patched.")
	}
	if user.Address.Address != userToPatch.Address.Address {
		t.Errorf("Email was not patched.")
	}
	if user.Updated.Sub(userToPatch.Updated) == 0 {
		t.Errorf("Last Updated not patched.")
	}
}

func TestDeleteUser(t *testing.T) {
	user := FindUserByID(users[0].ID)
	err := DeleteUser(user.ID)
	if err != nil {
		t.Error(err)
	}
	if users[0].ID == user.ID {
		t.Errorf("User was not deleted.")
	}
}
