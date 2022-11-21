package main

import (
	"fmt"
	"net/http"
	"testing"
)

func setup() {
	a := users["bappi"]
	a.Username = "bappi"
	a.password = "test"
	a.Name = "bharath"
	a.following = make(map[string]struct{})
	users["bappi"] = a
}

func TestUserExists(t *testing.T) {
	setup()
	username := "bappi"
	if !userExists(username) {
		t.Error("Something went wrong")
	}
}

func TestAddNewUser(t *testing.T) {
	a := "addusertest1"
	b := "password"
	c := "name"
	if !addNewUser(a, b, c) {
		t.Error("Something went wrong")
	}
}

func TestAuthentication(t *testing.T) {
	setup()
	username := "bappi"
	password := "test"

	if !authenticate(username, password) {
		t.Error("Test Authentication Failed")
	}

	invalidUsername := "syed"
	invalidPassword := "syed123"
	if authenticate(invalidUsername, invalidPassword) {
		t.Error("Test invalid Authentication Failed")
	}
}

func TestAddNewTweet(t *testing.T) {
	tweet := "test"
	setup()
	loggedInUser = "bappi"
	f := 0

	addNewTweet(tweet)
	temp := getUserTweets()
	for _, v := range temp.posts {
		if v == tweet {
			f = 1
			break
		}
	}

	if f == 0 {
		t.Error("Test Add new Tweet Failed")
	}

}

func Test(t *testing.T) {
	tweet := "test"
	setup()
	loggedInUser = "bappi"
	f := 0

	addNewTweet(tweet)
	temp := getUserTweets()
	for _, v := range temp.posts {
		if v == tweet {
			f = 1
			break
		}
	}

	if f == 0 {
		t.Error("Test Add new Tweet Failed")
	}

}


