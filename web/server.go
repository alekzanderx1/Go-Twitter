package main

import (
	"Twitter/authentication"
	"Twitter/tweets"
	"Twitter/users"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Failed to Listen to TCP: %v", err)
	}

	u := users.Server{}
	t := tweets.Server{}
	s := authentication.Server{}

	Server := grpc.NewServer()

	users.RegisterUserServiceServer(Server, &u)
	tweets.RegisterTweetsServiceServer(Server, &t)
	authentication.RegisterAuthServiceServer(Server, &s)

	if err := Server.Serve(lis); err != nil {
		log.Fatalf("Failed to Listen to TCP: %s", err)
	}
}
