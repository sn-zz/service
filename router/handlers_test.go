// Package router contains endpoint information for the service.
//
// sn - https://github.com/sn
package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sn/service/session"
	"github.com/sn/service/user"
)

var (
	server *httptest.Server
	client *http.Client
)

func TestIndex(t *testing.T) {
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Invalid response status code.")
	}
	if string(body) != "Welcome!\n" {
		t.Error("Invalid response body.")
	}

	users := user.GetAll()
	user := user.User{ID: users[0].ID, Password: "1@E4s67890"}
	authToken, err := getAuthToken(user)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", server.URL+"/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", string(authToken))
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("Invalid response status code.")
	}
	if string(body) != "Welcome, "+users[0].Username+"!\n" {
		t.Error("Invalid response body.")
	}
}

func TestMain(m *testing.M) {
	router := NewRouter()

	server = httptest.NewServer(router)
	client = &http.Client{}

	usernames := [4]string{"alex", "blake", "corey", "devon"}
	for _, un := range usernames {
		addr, err := mail.ParseAddress(strings.Title(un) + "<" + un + "@example.com>")
		if err != nil {
			log.Fatal(err)
		}
		u := user.User{Username: un, Password: "1@E4s67890", Address: addr, Created: time.Now()}
		u = user.Create(u)
		session.Create(u.ID)
	}

	os.Exit(m.Run())
}

func getAuthToken(user user.User) (string, error) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(user)
	resp, err := http.Post(server.URL+"/auth", "application/json; charset=utf-8", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	authToken, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Incorrect response status code.")
	}
	if len(authToken) == 0 {
		return "", fmt.Errorf("No authorization token was provided.")
	}
	return string(authToken), nil
}
