package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLoginRequestHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/loginre", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Error("Something went wrong")
	}

}

//signupRequestHandler

func TestSignupRequestHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/signupre", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Error("Something went wrong")
	}
	fmt.Println(req)
}

func TestAuthenticate(t *testing.T) {
	adddummy()
	username := "bappi"
	password := "test"

	if !authenticate(username, password) {
		t.Error("Something went wrong")
	}
}

func TestUserExists(t *testing.T) {
	adddummy()
	username := "bappi"
	if !UserExists(username) {
		t.Error("Something went wrong")
	}
}

func TestAddNewUser(t *testing.T) {
	a := "addusertest1"
	b := "password"
	c := "name"
	if !AddNewUser(a, b, c) {
		t.Error("Something went wrong")
	}

}
func TestAddNewTweet(t *testing.T) {
	tweet := "test"
	adddummy()
	loggedInUser = "bappi"
	f := 0

	AddNewTweet(tweet)
	//tweet = "test2"
	temp := users[loggedInUser]
	for _, v := range temp.posts {
		if v == tweet {
			f = 1
			break
		}
	}
	if f == 0 {
		t.Error("Something went wrong")
	}

}
