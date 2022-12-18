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

type Configuration struct {
	RaftClients []string
}

// Static variables
var CONFIG Configuration

// Load configuration from external file
func loadConfiguration() Configuration {
	file, err1 := os.Open("./tweets/tweet_config.json")
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

func (s *Server) GetTweetsByUsers(ctx context.Context, in *GetTweetsRequest) (*GetTweetsResponse, error) {
	var texts []string
	var createdBy []string
	var timestamps []string
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/tweets")
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
	raftUrl := findWorkingRAFTClient()

	resp, err := http.Get(raftUrl+"/tweets")
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
	cmd := exec.Command("curl", "-L", raftUrl+"tweets", "-XPUT", "-d "+string(dataBytes))

	cmd.Run()
	return &AddTweetResponse{Success: true}, nil
}
