package authentication

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Username  string
	Name      string
	Password  string
	Following map[string]struct{}
}

type Server struct {
	AuthServiceServer
}

// In memory non-persistent storage, to be replaced with database later
var data = make(map[string]User)

func (s *Server) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &data)
	temp := data[in.Username]
	result1 := temp.Password == in.Password
	if result1 {
		return &AuthenticateResponse{Success: true}, nil
	} else {

		return &AuthenticateResponse{Success: false}, nil
	}

}
