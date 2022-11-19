package main

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
)

var users = make(map[string]User)

type User struct {
	Username  string
	password  string
	followers [3]string
	posts     [5]string
}

var loggeduser string
var tp1 *template.Template

//users := make(map[string]int)

func signuprequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	temp := users[username]

	temp.Username = username
	temp.password = password
	users[username] = temp
	fmt.Println(users)
	http.ServeFile(res, req, "./static/login.html")
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/signup.html")

}
func loginrequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	if _, ok := users[username]; ok {
		fmt.Println(users[username])
		fmt.Println(password)
		temp := users[username]
		result1 := temp.password == password
		if result1 {

			fmt.Println("Login Success")
			loggeduser = username

			tp1.ExecuteTemplate(res, "welcome.html", loggeduser)

		} else {
			http.ServeFile(res, req, "./static/login.html")
			fmt.Println("Incorrect password")

		}

	} else {
		fmt.Println("No user")

		http.ServeFile(res, req, "./static/signup.html")

	}

}

func loginPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/login.html")

}
func userslist(res http.ResponseWriter, req *http.Request) {

	//var name = "bappi"
	fmt.Println(users)
	keys := reflect.ValueOf(users).MapKeys()
	type temp struct{
		people [] string
	}
	peeps:=temp(
		people:=keys,
	)
	tp1.ExecuteTemplate(res, "usersfeed.html", temp)

}

func main() {
	tp1, _ = tp1.ParseGlob("static/*.html")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signupre", signuprequest)
	http.HandleFunc("/loginre", loginrequest)
	http.HandleFunc("/users", userslist)

	http.HandleFunc("/login", loginPage)
	http.ListenAndServe(":8080", nil)

}
