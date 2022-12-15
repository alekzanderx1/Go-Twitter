package authentication

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	Username string
	Expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s *Server) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	var data = make(map[string]User)

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

	// Getting Sessions from raft
	resp_sessions, err := http.Get("http://127.0.0.1:12380/session")
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

		dataBytes, err := json.Marshal(sessions)
		if err != nil {
			fmt.Println(err)
		}
		cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/session", "-XPUT", "-d "+string(dataBytes))

		cmd.Run()
		return &AuthenticateResponse{Success: true, SessionToken: sessionToken}, nil
	} else {

		return &AuthenticateResponse{Success: false}, nil
	}

}

func (s *Server) ValidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
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
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, in.SessionToken)
		return &ValidateSessionResponse{Success: false}, nil
	}

	return &ValidateSessionResponse{Success: true, Username: userSession.Username}, nil
}

func (s *Server) InvalidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	resp_sessions, err := http.Get("http://127.0.0.1:12380/session")
	if err != nil {
		fmt.Println(err)
	}
	body_sessions, err := ioutil.ReadAll(resp_sessions.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(body_sessions, &sessions)
	delete(sessions, in.SessionToken)
	dataBytes, err := json.Marshal(sessions)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/session", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &ValidateSessionResponse{Success: true}, nil
}
