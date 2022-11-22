package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Definition of Structs for Data storage

type User struct {
	Username  string
	Name      string
	password  string
	following map[string]struct{}
	posts     []string
}

type Tweet struct {
	Text     string
	Username string
}

// In memory non-persistent storage

var users = make(map[string]User)
var loggedInUser = "Guest"
var tp1 *template.Template

// Authentication Methods

func userExists(username string) bool {
	if _, exists := users[username]; exists {
		return true
	} else {
		return false
	}

}

func authenticate(username string, password string) bool {
	temp := users[username]
	result1 := temp.password == password
	if result1 {
		return true
	} else {

		return false
	}

}

func addNewUser(username string, password string, name string) bool {

	temp := users[username]
	temp.Username = username
	temp.password = password
	temp.Name = name
	temp.following = make(map[string]struct{})
	users[username] = temp
	if userExists(username) {
		return true
	} else {
		return false
	}
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/signup.html")
}

func signupRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")

	if !userExists(username) {
		password := req.FormValue("password")
		name := req.FormValue("name")
		if addNewUser(username, password, name) {
			http.ServeFile(res, req, "./static/login.html")
		} else {
			http.Error(res, "Something went wrong", http.StatusConflict)
		}

	} else {
		http.Error(res, "Username Exists", http.StatusConflict)
	}

}

func loginPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/login.html")
}

func loginRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	if username == "" {
		http.Error(res, "Wrong Format for username", http.StatusFailedDependency)
		return

	}
	if userExists(username) {
		if authenticate(username, password) {
			fmt.Println("Login Success")
			loggedInUser = username
			userFeedHandler(res, req)
		} else {

			http.ServeFile(res, req, "./static/login.html")
			fmt.Println("Incorrect password")
		}

	} else {
		// User doesn't exist, prompt Signup
		http.ServeFile(res, req, "./static/signup.html")
	}
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser = "Guest"
	http.ServeFile(res, req, "./static/index.html")
}

// User Timeline Methods

func getTimelineForUser(user string) []Tweet {
	var tweets []Tweet

	following := users[user].following
	for friend, _ := range following {
		for _, tweet := range users[friend].posts {
			tweets = append(tweets, Tweet{Text: tweet, Username: friend})
		}
	}

	return tweets
}

func userFeedHandler(res http.ResponseWriter, req *http.Request) {
	type TemplateData struct {
		Username string
		Tweets   []Tweet
	}

	tweets := getTimelineForUser(loggedInUser)
	data := TemplateData{Username: loggedInUser, Tweets: tweets}
	tp1.ExecuteTemplate(res, "userfeed.html", data)
}

// User Follower Methods

func followUser(username string) {
	users[loggedInUser].following[username] = struct{}{}
}

func unfollowUser(username string) {
	delete(users[loggedInUser].following, username)
}

func getUserFollowers() map[string]struct{} {
	return users[loggedInUser].following
}

func followUserHandler(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	followUser(username)
	usersListHandler(res, req)
}

func unfollowUserHandler(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	unfollowUser(username)
	usersListHandler(res, req)
}

func usersListHandler(res http.ResponseWriter, req *http.Request) {
	type UserListItem struct {
		Username string
		Name     string
	}

	type TemplateData struct {
		FollowingList []UserListItem
		FollowList    []UserListItem
	}

	var followingUserList []UserListItem
	var followUserList []UserListItem

	following := getUserFollowers()

	for user, details := range users {
		if user != loggedInUser {
			_, follows := following[user]
			if follows == true {
				followingUserList = append(followingUserList, UserListItem{Username: user, Name: details.Name})
			} else {
				followUserList = append(followUserList, UserListItem{Username: user, Name: details.Name})
			}
		}
	}

	data := TemplateData{FollowingList: followingUserList, FollowList: followUserList}
	tp1.ExecuteTemplate(res, "users.html", &data)
}

// User Tweet Methods

func addNewTweet(tweet string) {
	temp := users[loggedInUser]
	temp.posts = append(temp.posts, tweet)
	users[loggedInUser] = temp
}

func getUserTweets() []string {
	return users[loggedInUser].posts
}

func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	tweet := req.FormValue("tweet")
	addNewTweet(tweet)
	tp1.ExecuteTemplate(res, "MyTweets.html", users[loggedInUser].posts)

}

func myTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	tp1.ExecuteTemplate(res, "MyTweets.html", getUserTweets())
}

func main() {
	tp1, _ = tp1.ParseGlob("static/*.html")
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signupre", signupRequestHandler)

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/loginre", loginRequestHandler)

	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/feed", userFeedHandler)

	http.HandleFunc("/users", usersListHandler)
	http.HandleFunc("/follow", followUserHandler)
	http.HandleFunc("/unfollow", unfollowUserHandler)

	http.HandleFunc("/tweet", newTweetRequestHandler)
	http.HandleFunc("/mytweets", myTweetRequestHandler)

	http.ListenAndServe(":8080", nil)
}
