package main

import (
	"Twitter/authentication"
	"Twitter/tweets"
	"Twitter/users"
	"time"

	"context"
	"html/template"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

// Definition of Structs for Data Conversion
type User struct {
	Username  string
	Name      string
	password  string
	following map[string]struct{}
	posts     []string
}

type Tweet struct {
	Text      string
	Username  string
	Timestamp string
}

// var loggedInUser = "Guest"
var tp1 *template.Template

// Helper method to validate session token for a incoming request
func validateSession(res http.ResponseWriter, req *http.Request) string {
	// We can obtain the session token from the requests cookies, which come with every request
	cookie, err1 := req.Cookie("session_token")
	if err1 != nil || cookie == nil {
		// If token not present, serve login/signup page
		http.ServeFile(res, req, "./static/index.html")
	}
	sessionToken := cookie.Value

	var conn *grpc.ClientConn
	conn, err2 := grpc.Dial(":9000", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("Couldn't connect: %s", err2)
	}

	a := authentication.NewAuthServiceClient(conn)

	response, err := a.ValidateSession(context.Background(), &authentication.ValidateSessionRequest{
		SessionToken: sessionToken,
	})

	if err != nil {
		log.Fatalf("Error when calling ValidateSession: %s", err)
	}

	// If session expired or token not valid, serve login/signup page
	if response.Success == false {
		http.ServeFile(res, req, "./static/index.html")
	}

	return response.Username
}

// Authentication Handlers

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
		Username: username,
		Password: password,
		Name:     name,
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

	a := authentication.NewAuthServiceClient(conn)

	response, err := a.Authenticate(context.Background(), &authentication.AuthenticateRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Error when calling Authenticate: %s", err)
		http.Error(res, "Something went wrong, please try again!", http.StatusConflict)
	}

	if !response.Success {
		http.Error(res, "Authentication Failed, please check username or password!", http.StatusConflict)
	} else {
		http.SetCookie(res, &http.Cookie{
			Name:    "session_token",
			Value:   response.SessionToken,
			Expires: time.Now().Add(120 * time.Second),
		})
		http.Redirect(res, req, "/feed", http.StatusSeeOther)
	}
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	validateSession(res, req)
	// We can obtain the session token from the requests cookies, which come with every request
	cookie, err1 := req.Cookie("session_token")
	if err1 != nil {
		// If token not present, serve login/signup page
		http.SetCookie(res, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	sessionToken := cookie.Value

	var conn *grpc.ClientConn
	conn, err2 := grpc.Dial(":9000", grpc.WithInsecure())
	if err2 != nil {
		log.Fatalf("Couldn't connect: %s", err2)
	}

	a := authentication.NewAuthServiceClient(conn)

	response, err := a.InvalidateSession(context.Background(), &authentication.ValidateSessionRequest{
		SessionToken: sessionToken,
	})

	if err != nil {
		log.Fatalf("Error when calling InvalidateSession: %s", err)
	}

	// If session expired or token not valid, serve login/signup page
	if response.Success == false {
		http.ServeFile(res, req, "./static/index.html")
	}

	http.SetCookie(res, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// User Timeline Handlers

func getTimelineForUser(user string, res http.ResponseWriter, req *http.Request) []Tweet {
	following := getUserFollowers(user, res, req).Following
	return getTweetsForUsers(following, res, req)
}

func userFeedHandler(res http.ResponseWriter, req *http.Request) {
	type TemplateData struct {
		Username string
		Tweets   []Tweet
	}
	loggedInUser := validateSession(res, req)
	tweets := getTimelineForUser(loggedInUser, res, req)
	data := TemplateData{Username: loggedInUser, Tweets: tweets}
	tp1.ExecuteTemplate(res, "userfeed.html", data)
}

// User Get, Follow, Unfollow Handlers

func getUserFollowers(username string, res http.ResponseWriter, req *http.Request) *users.GetFollowingResponse {
	loggedInUser := validateSession(res, req)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	u := users.NewUserServiceClient(conn)

	response, err := u.GetFollowers(context.Background(), &users.GetFollowingRequest{
		Username: loggedInUser,
	})

	if err != nil {
		log.Fatalf("Error when calling GetFollowers: %s", err)
		http.Error(res, "Something went wrong, please try again!", http.StatusConflict)
	}

	return response
}

func followUser(loggedInUser string, follow string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	u := users.NewUserServiceClient(conn)

	response, err := u.FollowUser(context.Background(), &users.AddFollowerRequest{
		Username: loggedInUser,
		Follow:   follow,
	})

	if err != nil || !response.Success {
		log.Fatalf("Error when calling FollowUser: %s", err)
	}
}

func unfollowUser(loggedInUser string, follow string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	u := users.NewUserServiceClient(conn)

	response, err := u.UnfollowUser(context.Background(), &users.RemoveFollowerRequest{
		Username: loggedInUser,
		Follow:   follow,
	})

	if err != nil || !response.Success {
		log.Fatalf("Error when calling UnfollowUser: %s", err)
	}
}

func followUserHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser := validateSession(res, req)
	username := req.FormValue("username")
	followUser(loggedInUser, username)
	usersListHandler(res, req)
}

func unfollowUserHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser := validateSession(res, req)
	username := req.FormValue("username")
	unfollowUser(loggedInUser, username)
	usersListHandler(res, req)
}

func usersListHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser := validateSession(res, req)
	type TemplateData struct {
		FollowingList []string
		FollowList    []string
	}
	response := getUserFollowers(loggedInUser, res, req)
	data := TemplateData{FollowingList: response.Following, FollowList: response.Suggestions}
	tp1.ExecuteTemplate(res, "users.html", &data)
}

// User Tweet Handlers

func getTweetsForUsers(usernames []string, res http.ResponseWriter, req *http.Request) []Tweet {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	t := tweets.NewTweetsServiceClient(conn)

	response, err := t.GetTweetsByUsers(context.Background(), &tweets.GetTweetsRequest{
		Usernames: usernames,
	})

	if err != nil {
		log.Fatalf("Error when calling GetTweetsByUsers: %s", err)
	}

	var result []Tweet
	for i := 0; i < len(response.Text); i++ {
		tweet := Tweet{
			Text:      response.Text[i],
			Username:  response.CreatedBy[i],
			Timestamp: response.Timestamp[i],
		}
		result = append(result, tweet)
	}

	return result
}

func addNewTweet(tweet string, username string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %s", err)
	}

	t := tweets.NewTweetsServiceClient(conn)

	response, err := t.AddNewTweet(context.Background(), &tweets.AddTweetRequest{
		Username: username,
		Text:     tweet,
	})

	if err != nil || !response.Success {
		log.Fatalf("Error when calling AddNewTweet: %s", err)
	}
}

func newTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser := validateSession(res, req)
	tweet := req.FormValue("tweet")
	addNewTweet(tweet, loggedInUser)
	myTweetRequestHandler(res, req)
}

func myTweetRequestHandler(res http.ResponseWriter, req *http.Request) {
	loggedInUser := validateSession(res, req)
	type TemplateData struct {
		Tweets []Tweet
	}
	usernames := []string{loggedInUser}
	userTweets := getTweetsForUsers(usernames, res, req)
	data := TemplateData{Tweets: userTweets}
	tp1.ExecuteTemplate(res, "MyTweets.html", data)
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
