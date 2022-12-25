package authentication

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterAuthServiceServer(s, &Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAuthenticationEndToEnd(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := NewAuthServiceClient(conn)

	// Authenticate a User whom we added during Users service test
	resp, err := client.Authenticate(ctx, &AuthenticateRequest{Username: "test1", Password: "password1"})
	if err != nil {
		t.Fatalf("Test Authenticate failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	if !resp.Success {
		t.Error("Test Authenticate Failed")
	}
	sessionToken := resp.SessionToken

	time.Sleep(1 * time.Second)

	// Now validate the session
	resp2, err2 := client.ValidateSession(ctx, &ValidateSessionRequest{SessionToken: sessionToken})
	if err2 != nil {
		t.Fatalf("Test ValidateSession failed: %v", err2)
	}
	log.Printf("Response: %+v", resp2)

	if !resp2.Success || resp2.Username != "test1" {
		t.Error("Test ValidateSession Failed")
	}

	// Now Invalidate the session
	resp3, err3 := client.InvalidateSession(ctx, &ValidateSessionRequest{SessionToken: sessionToken})
	if err3 != nil {
		t.Fatalf("Test InvalidateSession failed: %v", err3)
	}
	log.Printf("Response: %+v", resp3)

	if !resp3.Success {
		t.Error("Test InvalidateSession Failed")
	}

	time.Sleep(1 * time.Second)

	// Now try validation with invalidated token
	resp4, err4 := client.ValidateSession(ctx, &ValidateSessionRequest{SessionToken: sessionToken})
	if err4 != nil {
		t.Fatalf("Test ValidateSession after logout Failed: %v", err4)
	}
	log.Printf("Response: %+v", resp4)

	if resp4.Success {
		t.Error("Test ValidateSession after logout Failed")
	}
}
