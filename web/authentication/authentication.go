package authentication

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	username string
	expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
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

	if result1 {
		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(120 * time.Second)

		// Set the token in the session map, along with the session information
		sessions[sessionToken] = session{
			username: in.Username,
			expiry:   expiresAt,
		}

		return &AuthenticateResponse{Success: true, SessionToken: sessionToken}, nil
	} else {

		return &AuthenticateResponse{Success: false}, nil
	}

}

func (s *Server) ValidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	// Get the session from our session map
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

	return &ValidateSessionResponse{Success: true, Username: userSession.username}, nil
}

func (s *Server) InvalidateSession(ctx context.Context, in *ValidateSessionRequest) (*ValidateSessionResponse, error) {
	delete(sessions, in.SessionToken)
	return &ValidateSessionResponse{Success: true}, nil
}
