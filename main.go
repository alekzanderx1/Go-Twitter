package main



import (
	
	"net/http"
	"fmt"
	
)
var users = make(map[string]string)
//users := make(map[string]int)

func signuprequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	users[username]=password
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
	if _, ok := users[username]; ok{
		fmt.Println(users[username])
		fmt.Println(password)
		result1 := users[username] == password
		if result1{
			http.ServeFile(res, req, "./static/welcome.html")
			fmt.Println("Login Success")
		}else{
			http.ServeFile(res, req, "./static/login.html")
			fmt.Println("Incorrect password")


		}

	}else{
		fmt.Println("No user")

	http.ServeFile(res, req, "./static/signup.html")

	}


}

func loginPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/login.html")

}





func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}
func public() http.Handler {
	return http.StripPrefix("/", http.FileServer(http.Dir("./indext.html")))
}



func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signupre", signuprequest)
	http.HandleFunc("/loginre", loginrequest)
	


	http.HandleFunc("/login", loginPage)
	http.ListenAndServe(":8080", nil)
	
}
