package main



import (
	"fmt"
	
	"log"
	"net/http"
	"os"
	"time"
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

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		log.Println(req)
		next.ServeHTTP(w, r)
		log.Println(req, "completed in", time.Now().Sub(start))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", logging(public()))
	//mux.Handle("/", logging(index()))
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Println("main: running simple server on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: couldn't start simple server: %v\n", err)
	}
	
}
