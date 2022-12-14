package tweets

import (
	context "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

// Definition of Structs for Data storage
type Tweet struct {
	Text             string
	CreatedBy        string
	CreatedTimestamp string
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
	resp, err := http.Get("http://127.0.0.1:12380/tweets")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &tweets)
	for _, user := range in.Usernames {
		user_tweets := tweets[user]
		for _, user_tweet := range user_tweets {
			texts = append(texts, user_tweet.Text)
			createdBy = append(createdBy, user_tweet.CreatedBy)
			timestamps = append(timestamps, user_tweet.CreatedTimestamp)
		}
	}

	return &GetTweetsResponse{Text: texts, CreatedBy: createdBy, Timestamp: timestamps}, nil
}

func (s *Server) AddNewTweet(ctx context.Context, in *AddTweetRequest) (*AddTweetResponse, error) {
	resp, err := http.Get("http://127.0.0.1:12380/tweets")
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal(body, &tweets)

	temp := Tweet{Text: in.Text, CreatedBy: in.Username, CreatedTimestamp: time.Now().Format("2006-01-02 15:04:05")}
	tweets[in.Username] = append(tweets[in.Username], temp)

	dataBytes, err := json.Marshal(tweets)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("curl", "-L", "http://127.0.0.1:12380/tweets", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &AddTweetResponse{Success: true}, nil
}
