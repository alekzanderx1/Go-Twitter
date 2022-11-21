package main

import (
	"fmt"
	"html/template"
	"net/http"
	// unable to import time package @syed take a look
)

var users = make(map[string]User)

type User struct {
	Username  string
	Name      string
	password  string
	following map[string]struct{}
	posts     []string
}

var loggedInUser = "Guest"
var tp1 *template.Template

func adddummy() {
	a := users["bappi"]
	a.Username = "bappi"
	a.password = "test"
	a.Name = "bharath"
	a.following = make(map[string]struct{})
	users["bappi"] = a

}
func UserExists(username string) bool {
	//adding dummy for unit test cases.

	if _, exists := users[username]; exists {
		return true
	} else {
		return false
	}

}
func AddNewUser(username string, password string, name string) bool {

	temp := users[username]
	temp.Username = username
	temp.password = password
	temp.Name = name
	temp.following = make(map[string]struct{})
	users[username] = temp
	if UserExists(username) {
		return true
	} else {
		return false
	}
}
func signupRequestHandler(res http.ResponseWriter, req *http.Request) {
	//fmt.Println(req)
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")

	if !UserExists(username) {
		password := req.FormValue("password")
		name := req.FormValue("name")
		if AddNewUser(username, password, name) {
			http.ServeFile(res, req, "./static/login.html")
		} else {
			http.Error(res, "Something went wrong", http.StatusConflict)
		}

	} else {
		http.Error(res, "Username Exists", http.StatusConflict)
	}

}

func signupPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/signup.html")
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
	if UserExists(username) {
		fmt.Println(users[username])
		fmt.Println(password)

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

func loginPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/login.html")

}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser = "Guest"
	http.ServeFile(res, req, "./static/index.html")
}

func userFeedHandler(res http.ResponseWriter, req *http.Request) {
	type Tweet struct {
		Text     string
		Username string
	}

	type Data struct {
		Username string
		Tweets   []Tweet
	}

	var tweets []Tweet

	following := users[loggedInUser].following

	for friend, _ := range following {
		for _, tweet := range users[friend].posts {
			tweets = append(tweets, Tweet{Text: tweet, Username: friend})
		}
	}

	data := Data{Username: loggedInUser, Tweets: tweets}
	tp1.ExecuteTemplate(res, "userfeed.html", data)
}

func usersListHandler(res http.ResponseWriter, req *http.Request) {
	type UserListItem struct {
		Username string
		Name     string
	}

	type Data struct {
		FollowingList []UserListItem
		FollowList    []UserListItem
	}

	var followingUserList []UserListItem
	var followUserList []UserListItem

	following := users[loggedInUser].following

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

	data := Data{FollowingList: followingUserList, FollowList: followUserList}
	tp1.ExecuteTemplate(res, "users.html", &data)
}

func AddNewTweet(tweet string) {
	//tweet := req.FormValue("tweet")
	temp := users[loggedInUser]
	temp.posts = append(temp.posts, tweet)
	users[loggedInUser] = temp

}

func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	// Throw error for Guest
	tweet := req.FormValue("tweet")
	AddNewTweet(tweet)
	fmt.Println(users)
	tp1.ExecuteTemplate(res, "MyTweets.html", users[loggedInUser].posts)

}

func myTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	// Redirect to Login for Guest
	tp1.ExecuteTemplate(res, "MyTweets.html", users[loggedInUser].posts)
}

func followUserHandler(res http.ResponseWriter, req *http.Request) {
	// Throw error for Guest and Not POST
	username := req.FormValue("username")
	users[loggedInUser].following[username] = struct{}{}
	usersListHandler(res, req)
}

func unfollowUserHandler(res http.ResponseWriter, req *http.Request) {
	//Throw error for Guest and Not POST
	username := req.FormValue("username")
	delete(users[loggedInUser].following, username)
	usersListHandler(res, req)
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

	//http.ListenAndServe("0.0.0.0:8000", nil)
	http.ListenAndServe(":8080", nil)
}
