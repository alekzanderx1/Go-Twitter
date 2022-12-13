package main

import (
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
		t.Error("User Exists Test Went Wrong")
	}
	username = "bappi_2"
	if userExists(username) {
		t.Error("User Exists Test Went Wrong")

	}
}

func TestAddNewUser(t *testing.T) {
	a := "addusertest1"
	b := "password"
	c := "name"
	if !addNewUser(a, b, c) {
		t.Error("Test Add New User Failed")
	}
	setup()

	//trying to add existing user
	a = "bappi"
	b = "test"
	c = "bharath"
	if addNewUser(a, b, c) {
		t.Error("Test Add New User Failed, User Already Exists")
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

	//valid usernme and invalid password
	if authenticate(username, invalidPassword) {
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
	for _, v := range temp {
		if v == tweet {
			f = 1
			break
		}
	}

	if f == 0 {
		t.Error("Test Add new Tweet Failed")

	}
	//Testing with tweet which is not tweet
	f = 0
	newtweet := "test_2"
	for _, v := range temp {
		if v == newtweet {
			f = 1
			break
		}
	}
	if f == 1 {
		t.Error("Unknown Tweet found")

	}

}
