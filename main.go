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
	password  string
	following []string
	posts     []string
}

type UserListItem struct {
	Username  string
	following bool
}

var loggedInUser string
var tp1 *template.Template

func signupRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	temp := users[username]

	temp.Username = username
	temp.password = password
	users[username] = temp
	http.ServeFile(res, req, "./static/login.html")
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/signup.html")

}
func loginRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	if _, exists := users[username]; exists {
		fmt.Println(users[username])
		fmt.Println(password)
		temp := users[username]
		result1 := temp.password == password
		if result1 {
			fmt.Println("Login Success")
			loggedInUser = username
			tp1.ExecuteTemplate(res, "userfeed.html", loggedInUser)
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

func userFeedHandler(res http.ResponseWriter, req *http.Request) {
	tp1.ExecuteTemplate(res, "userfeed.html", loggedInUser)

}

func usersListHandler(res http.ResponseWriter, req *http.Request) {

	var userList []UserListItem
	for user, _ := range users {
		if user != loggedInUser {
			userList = append(userList, UserListItem{Username: user})
		}
	}

	tp1.ExecuteTemplate(res, "users.html", userList)

}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser = ""
	http.ServeFile(res, req, "./static/index.html")
}

func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	tweet := req.FormValue("tweet")
	temp := users[loggedInUser]
	temp.posts = append(temp.posts, tweet)
	users[loggedInUser] = temp
	fmt.Println(users)
	tp1.ExecuteTemplate(res, "MyTweets.html", users[loggedInUser].posts)

}

// myTweetRequestHandler
func myTweetRequestHandler(res http.ResponseWriter, req *http.Request) {

	tp1.ExecuteTemplate(res, "MyTweets.html", users[loggedInUser].posts)

}

func main() {
	tp1, _ = tp1.ParseGlob("static/*.html")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/feed", userFeedHandler)

	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/users", usersListHandler)
	http.HandleFunc("/signupre", signupRequestHandler)
	http.HandleFunc("/loginre", loginRequestHandler)
	http.HandleFunc("/tweet", newTweetRequestHandler)
	//mytweets
	http.HandleFunc("/mytweets", myTweetRequestHandler)

	http.ListenAndServe(":8080", nil)

}
