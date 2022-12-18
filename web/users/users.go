package users

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"os"
	"log"
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

type Configuration struct {
	RaftClients []string
}

// Static variables
var CONFIG Configuration

// Load configuration from external file
func loadConfiguration() Configuration {
	file, err1 := os.Open("./users/users_config.json")
	if err1 != nil {
		fmt.Print("File reading error")
		fmt.Print(err1)
	}
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	file.Close()
	return conf
}

func findWorkingRAFTClient() string {
	fmt.Print(CONFIG.RaftClients)
	for _, url := range CONFIG.RaftClients {
		_, err := http.Get(url + "/ping")
		if err == nil {
			return url
		}
	}
	log.Fatalf("Couldn't connect find working RAFT client")
	return ""
}

func init() {
	CONFIG = loadConfiguration()
}

func (s *Server) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/users")
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
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/users")
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
	cmd := exec.Command("curl", "-L", raftUrl+"/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()

	return &AddUserResponse{Success: true}, nil
}

func (s *Server) GetFollowers(ctx context.Context, in *GetFollowingRequest) (*GetFollowingResponse, error) {
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/users")
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
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/users")
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
	cmd := exec.Command("curl", "-L", raftUrl+"/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &AddFollowerResponse{Success: true}, nil
}

func (s *Server) UnfollowUser(ctx context.Context, in *RemoveFollowerRequest) (*RemoveFollowerResponse, error) {
	raftUrl := findWorkingRAFTClient()
	
	resp, err := http.Get(raftUrl+"/users")
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
	cmd := exec.Command("curl", "-L", raftUrl+"/users", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &RemoveFollowerResponse{Success: true}, nil
}
