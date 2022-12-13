package main

import (
	"html/template"
	"net/http"
	"Twitter/users"
	"google.golang.org/grpc"
	"log"
	"context"
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

var data = make(map[string]User)
var loggedInUser = "Guest"
var tp1 *template.Template

// Authentication Methods
func signupPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/signup.html")
}

func signupRequestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	name := req.FormValue("name")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	u := users.NewUserServiceClient(conn)

	response, err := u.AddNewUser(context.Background(), &users.AddUserRequest{
		Username:        username,
		Password:        password,
		Name:            name,
	})

	if err != nil {
		log.Fatalf("Error when calling AddNewUser: %s", err)
	}

	if !response.Success {
		http.Error(res, "Something went wrong, make sure user doesn't exist already!", http.StatusConflict)
	} else {
		http.ServeFile(res, req, "./static/login.html")
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

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	u := users.NewUserServiceClient(conn)

	response, err := u.Authenticate(context.Background(), &users.AuthenticateRequest{
		Username:        username,
		Password:        password,
	})

	if err != nil {
		log.Fatalf("Error when calling Authenticate: %s", err)
	}

	if !response.Success {
		http.Error(res, "Authentication Failed", http.StatusConflict)
		http.ServeFile(res, req, "./static/login.html")
	} else {
		loggedInUser = username
		userFeedHandler(res, req)
	}

	// if userExists(username) {
	// 	if authenticate(username, password) {
	// 		fmt.Println("Login Success")
	// 		loggedInUser = username
	// 		userFeedHandler(res, req)
	// 	} else {

	// 		http.ServeFile(res, req, "./static/login.html")
	// 		fmt.Println("Incorrect password")
	// 	}

	// } else {
	// 	// User doesn't exist, prompt Signup
	// 	http.ServeFile(res, req, "./static/signup.html")
	// }
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser = "Guest"
	http.ServeFile(res, req, "./static/index.html")
}

// User Timeline Methods

func getTimelineForUser(user string) []Tweet {
	var tweets []Tweet

	following := data[user].following
	for friend, _ := range following {
		for _, tweet := range data[friend].posts {
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
	data[loggedInUser].following[username] = struct{}{}
}

func unfollowUser(username string) {
	delete(data[loggedInUser].following, username)
}

func getUserFollowers() map[string]struct{} {
	return data[loggedInUser].following
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

	for user, details := range data {
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
	temp := data[loggedInUser]
	temp.posts = append(temp.posts, tweet)
	data[loggedInUser] = temp
}

func getUserTweets() []string {
	return data[loggedInUser].posts
}

func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	tweet := req.FormValue("tweet")
	addNewTweet(tweet)
	tp1.ExecuteTemplate(res, "MyTweets.html", data[loggedInUser].posts)

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

	http.ListenAndServe("0.0.0.0:8000", nil)
}
