package users

import (
	"testing"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/grpc"
	"log"
	"context"
	"net"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    RegisterUserServiceServer(s, &Server{})
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

func TestAddNewUser(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    client := NewUserServiceClient(conn)
    resp, err := client.AddNewUser(ctx, &AddUserRequest{Username: "test0", Password:"password0", Name:"RandomName0"})
    if err != nil {
        t.Fatalf("AddNewUser failed: %v", err)
    }
    log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Add New User Failed")
	}
}

func TestAddRemoveAndGetFollowers(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    client := NewUserServiceClient(conn)

	// Add a few users
    resp, err := client.AddNewUser(ctx, &AddUserRequest{Username: "test1", Password:"password1", Name:"RandomName1"})
    if err != nil {
        t.Fatalf("AddNewUser GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Add New User Failed")
	}

	resp, err = client.AddNewUser(ctx, &AddUserRequest{Username: "test2", Password:"password2", Name:"RandomName2"})
    if err != nil {
        t.Fatalf("AddNewUser GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Add New User Failed")
	}

	// Make user1 follow user2
	resp2, err := client.FollowUser(ctx, &AddFollowerRequest{Username: "test1", Follow:"test2"})
    if err != nil {
        t.Fatalf("FollowUser GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp2)

	if !resp2.Success {
		t.Error("Test FollowUser Failed")
	}

	// Check user1 followers contian user2
	resp3, err := client.GetFollowers(ctx, &GetFollowingRequest{Username: "test1"})
    if err != nil {
        t.Fatalf("GetFollowers GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp3)

	if !contains(resp3.Following, "test2") {
		t.Error("Test GetFollowers Failed")
	} 

	// Make user1 unfollow user2
	resp4, err := client.UnfollowUser(ctx, &RemoveFollowerRequest{Username: "test1", Follow:"test2"})
    if err != nil {
        t.Fatalf("UnfollowUser GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp4)

	if !resp4.Success {
		t.Error("Test UnfollowUser Failed")
	}

	// Check user1 followers doesn't contain user2
	resp5, err := client.GetFollowers(ctx, &GetFollowingRequest{Username: "test1"})
    if err != nil {
        t.Fatalf("FollowUser GRPC failed: %v", err)
    }
    log.Printf("Response: %+v", resp5)

	if contains(resp5.Following, "test2") {
		t.Error("Test GetFollowes Failed")
	} 
}

