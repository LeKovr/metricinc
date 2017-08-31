//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"

	pb "lekovr/exam/counter"
	"lekovr/exam/lib/server"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv, _ := server.NewServer()
	s := grpc.NewServer()
	pb.RegisterCounterServer(s, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
