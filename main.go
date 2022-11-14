package main



import (
	
	"net/http"
	
)


func signupPage(res http.ResponseWriter, req *http.Request) {
}


func loginPage(res http.ResponseWriter, req *http.Request) {
}





func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}
func public() http.Handler {
	return http.StripPrefix("/", http.FileServer(http.Dir("./indext.html")))
}



func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
	
}
