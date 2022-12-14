package tweets

import(
	context "context"
	"time"
)

// Definition of Structs for Data storage
type Tweet struct {
	text             string
	createdBy        string
	createdTimestamp string
}

type Server struct {
	TweetsServiceServer
}

// In memory non-persistent storage
var tweets = make(map[string][]Tweet)


func (s *Server) GetTweetsByUsers(ctx context.Context, in *GetTweetsRequest) (*GetTweetsResponse, error) {
	var texts []string
	var createdBy []string
	var timestamps []string

	for _, user := range in.Usernames {
		user_tweets := tweets[user]
		for _, user_tweet := range user_tweets {
			texts = append(texts, user_tweet.text)
			createdBy = append(createdBy, user_tweet.createdBy)
			timestamps = append(timestamps, user_tweet.createdTimestamp)
		}
	}

	return &GetTweetsResponse{Text: texts, CreatedBy: createdBy, Timestamp: timestamps}, nil
}

func (s *Server) AddNewTweet(ctx context.Context, in *AddTweetRequest) (*AddTweetResponse, error){
	temp := Tweet{text: in.Text, createdBy: in.Username, createdTimestamp:  time.Now().Format("2006-01-02 15:04:05") }
	tweets[in.Username] = append(tweets[in.Username],temp)
	return &AddTweetResponse{Success: true}, nil
}
