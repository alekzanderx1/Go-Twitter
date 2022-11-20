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
	Name 	  string
	password  string
	following map[string]struct{}
	posts     []string
}

var loggedInUser = "Guest"
var tp1 *template.Template

func signupRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	name := req.FormValue("name")
	temp := users[username]

	temp.Username = username
	temp.password = password
	temp.Name = name
	temp.following = make(map[string]struct{})
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

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser = "Guest"
	http.ServeFile(res, req, "./static/index.html")
}

func userFeedHandler(res http.ResponseWriter, req *http.Request) {
	tp1.ExecuteTemplate(res, "userfeed.html", loggedInUser)
}

func usersListHandler(res http.ResponseWriter, req *http.Request) {
	type UserListItem struct {
		Username  string
		Name string
	}

	type Data struct {
		FollowingList []UserListItem
		FollowList []UserListItem
	}

	var followingUserList []UserListItem
	var followUserList []UserListItem

	following := users[loggedInUser].following

	for user, details := range users {
		if user != loggedInUser {
			_, follows  := following[user]
			if follows == true {
				followingUserList = append(followingUserList, UserListItem{Username: user, Name: details.Name})
			} else {
				followUserList = append(followUserList, UserListItem{Username: user, Name: details.Name})
			}
		}
	}

	data := Data{FollowingList:followingUserList,FollowList:followUserList}
	tp1.ExecuteTemplate(res, "users.html", &data)
}


func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	// Throw error for Guest 
	tweet := req.FormValue("tweet")
	temp := users[loggedInUser]
	temp.posts = append(temp.posts, tweet)
	users[loggedInUser] = temp
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
	usersListHandler(res,req)
}

func unfollowUserHandler(res http.ResponseWriter, req *http.Request) {
	//Throw error for Guest and Not POST
	username := req.FormValue("username")
	delete(users[loggedInUser].following,username)
	usersListHandler(res,req)
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
