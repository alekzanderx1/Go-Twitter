package main

import (
	"fmt"
	"log"
	"net"
	"Twitter/users"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Failed to Listen to TCP: %v", err)
	}

	u := users.Server{}
	Server := grpc.NewServer()

	users.RegisterUserServiceServer(Server, &u)

	if err := Server.Serve(lis); err != nil {
		log.Fatalf("Failed to Listen to TCP: %s", err)
	}
}