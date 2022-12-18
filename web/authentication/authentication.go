package authentication

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
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

// each session contains the username of the user and the time at which it expires
type session struct {
	Username string
	Expiry   time.Time
}

type Configuration struct {
	RAFT_CLIENTS []string
}

// Static variables
var CONFIG Configuration

func init() {
	CONFIG = loadConfiguration()
}

// Load configuration from external file
func loadConfiguration() Configuration {
	file, _ := os.Open("conf.json")
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
	for _, url := range CONFIG.RAFT_CLIENTS {
		_, err := http.Get(url + "/ping")
		if err == nil {
			return url
		}
	}
	log.Fatalf("Couldn't connect find working RAFT client")
	return ""
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s *Server) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	var sessions = map[string]session{}
	var data = make(map[string]User)
	raftUrl := findWorkingRAFTClient()
	resp, err := http.Get(raftUrl + "/users")
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

	// Getting Sessions from raft
	resp_sessions, err := http.Get(raftUrl + "/session")
	if err != nil {
		fmt.Println(err)
	}
	body_sessions, err := ioutil.ReadAll(resp_sessions.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body_sessions, &sessions)

	if result1 {
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		// Set the token in the session map, along with the session information
		sessions[sessionToken] = session{
			Username: in.Username,
			Expiry:   expiresAt,
		}

		// Persist changes to Raft
		dataBytes, err := json.Marshal(sessions)
		if err != nil {
			fmt.Println(err)
		}
		cmd := exec.Command("curl", "-L", raftUrl+"/session", "-XPUT", "-d "+string(dataBytes))
		cmd.Run()
		time.Sleep(1 * time.Second)
		return &AuthenticateResponse{Success: true, SessionToken: sessionToken}, nil
	} else {

		return &AuthenticateResponse{Success: false}, nil
	}

}

func (s *Server) ValidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	var sessions = map[string]session{}
	// Get the session from our session map
	resp_sessions, err := http.Get("http://127.0.0.1:12380/session")
	if err != nil {
		fmt.Println(err)
	}
	body_sessions, err := ioutil.ReadAll(resp_sessions.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body_sessions, &sessions)

	userSession, exists := sessions[in.SessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		return &ValidateSessionResponse{Success: false}, nil
	}
	// If the session is present, but has expired, we can delete the session
	if userSession.isExpired() {
		delete(sessions, in.SessionToken)

		// Persist session changes to Raft
		dataBytes, err := json.Marshal(sessions)
		if err != nil {
			fmt.Println(err)
		}
		cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/session", "-XPUT", "-d "+string(dataBytes))
		cmd.Run()
		time.Sleep(1 * time.Second)
		return &ValidateSessionResponse{Success: false}, nil
	}

	return &ValidateSessionResponse{Success: true, Username: userSession.Username}, nil
}

func (s *Server) InvalidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	var sessions = map[string]session{}
	// Get session data from raft
	resp_sessions, err := http.Get("http://127.0.0.1:12380/session")
	if err != nil {
		fmt.Println(err)
	}
	body_sessions, err := ioutil.ReadAll(resp_sessions.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body_sessions, &sessions)

	// Delete token
	delete(sessions, in.SessionToken)

	// Persist changes to raft
	dataBytes, err := json.Marshal(sessions)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/session", "-XPUT", "-d "+string(dataBytes))
	cmd.Run()
	time.Sleep(1 * time.Second)
	return &ValidateSessionResponse{Success: true}, nil
}
