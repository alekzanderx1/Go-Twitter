package users

import(
	context "context"
)

// Definition of Structs for Data storage
type User struct {
	Username  string
	Name      string
	password  string
	following map[string]struct{}
}

type Server struct {
	UserServiceServer
}

// In memory non-persistent storage, to be replaced with database later
var data = make(map[string]User)



func (s *Server) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	temp := data[in.Username]
	result1 := temp.password == in.Password
	if result1 {
		return &AuthenticateResponse{Success: true}, nil
	} else {

		return &AuthenticateResponse{Success: false}, nil
	}

}

func (s *Server) AddNewUser(ctx context.Context, in *AddUserRequest)  (*AddUserResponse, error) {
	if _, exists := data[in.Username]; exists {
		return &AddUserResponse{Success: false}, nil
	}
	temp := data[in.Username]
	temp.Username = in.Username
	temp.password = in.Password
	temp.Name = in.Name
	temp.following = make(map[string]struct{})
	data[in.Username] = temp
	return &AddUserResponse{Success: true}, nil
}


func (s *Server)  GetFollowers(ctx context.Context, in *GetFollowingRequest)  (*GetFollowingResponse, error) {
	following := data[in.Username].following
	followingResponse := []string{}
	suggestions := []string{}

	for user := range following { 
		followingResponse = append(followingResponse,user)
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

func (s *Server)  FollowUser(ctx context.Context, in *AddFollowerRequest)  (*AddFollowerResponse, error) {
	data[in.Username].following[in.Follow] = struct{}{}
	return &AddFollowerResponse{Success: true}, nil
}

func (s *Server) UnfollowUser(ctx context.Context, in *RemoveFollowerRequest)  (*RemoveFollowerResponse, error) {
	delete(data[in.Username].following, in.Follow)
	return &RemoveFollowerResponse{Success: true}, nil
}
