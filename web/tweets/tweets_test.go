package tweets

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterTweetsServiceServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestAddNewTweet(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewTweetsServiceClient(conn)
	resp, err := client.AddNewTweet(ctx, &AddTweetRequest{Text: "Test tweet!", Username: "test1"})
	if err != nil {
		t.Fatalf("AddNewTweet failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Add New Tweet Failed")
	}
}

func TestAddAndRetriveTweets(t *testing.T) {
	time.Sleep(1 * time.Second)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewTweetsServiceClient(conn)

	// Add a new tweet
	resp, err := client.AddNewTweet(ctx, &AddTweetRequest{Text: "Test tweet2!", Username: "test1"})
	if err != nil {
		t.Fatalf("AddNewTweet failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Add New Tweet Failed")
	}

	time.Sleep(1 * time.Second)

	// Retrieve tweets for given user
	usernames := []string{"test1"}
	resp2, err := client.GetTweetsByUsers(ctx, &GetTweetsRequest{Usernames: usernames})
	if err != nil {
		t.Fatalf("GetTweetsByUsers failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	// There should be two tweets for given user now, one from first test and one uploaded above
	if len(resp2.Text) != 2 {
		t.Error("Test Get Tweets Failed")
	}
}
