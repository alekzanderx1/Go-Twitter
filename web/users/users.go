package users

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// Definition of Structs for Data storage
type User struct {
	Username  string
	Name      string
	Password  string
	Following map[string]struct{}
}

type Server struct {
	UserServiceServer
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

func (s *Server) AddNewUser(ctx context.Context, in *AddUserRequest) (*AddUserResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &data)
	if _, exists := data[in.Username]; exists {
		return &AddUserResponse{Success: false}, nil
	}

	temp := data[in.Username]
	temp.Username = in.Username
	temp.Password = in.Password
	temp.Name = in.Name
	temp.Following = make(map[string]struct{})
	data[in.Username] = temp

	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()

	return &AddUserResponse{Success: true}, nil
}

func (s *Server) GetFollowers(ctx context.Context, in *GetFollowingRequest) (*GetFollowingResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &data)
	following := data[in.Username].Following
	followingResponse := []string{}
	suggestions := []string{}

	for user := range following {
		followingResponse = append(followingResponse, user)
	}

	for user := range data {
		if user != in.Username {
			_, follows := following[user]
			if !follows {
				suggestions = append(suggestions, user)
			}
		}
	}

	return &GetFollowingResponse{Following: followingResponse, Suggestions: suggestions}, nil
}

func (s *Server) FollowUser(ctx context.Context, in *AddFollowerRequest) (*AddFollowerResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &data)
	data[in.Username].Following[in.Follow] = struct{}{}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &AddFollowerResponse{Success: true}, nil
}

func (s *Server) UnfollowUser(ctx context.Context, in *RemoveFollowerRequest) (*RemoveFollowerResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/users")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &data)
	delete(data[in.Username].Following, in.Follow)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &RemoveFollowerResponse{Success: true}, nil
}
