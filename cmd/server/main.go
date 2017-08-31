package main

import (
	"log"
	"net"

	//"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"
	"github.com/jessevdk/go-flags"
	pb "lekovr/exam/counter"
	"lekovr/exam/lib/boltdb"
	"lekovr/exam/lib/server"
)

type Config struct {
	Listen string `long:"listen" default:":50051" description:""`

	Server server.Config
	Store  boltdb.Config
}

const (
	port = ":50051"
	file = "data.db"
)

func main() {

	log.Println("Starting")

	var cfg Config
	_, err := flags.Parse(&cfg)
	StopOnError(err, "parse opts")
	log.Printf("Starting conf %+v", cfg)

	lis, err := net.Listen("tcp", port)
	StopOnError(err, "listen socket")
	log.Println("Starting")

	store, err := boltdb.Open(cfg.Store, file)
	srv, err := server.NewServer(cfg.Server, store)
	StopOnError(err, "db connect")
	defer srv.Close()
	log.Println("Starting")
	s := grpc.NewServer()
	pb.RegisterCounterServer(s, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StopOnError(e error, d string) {
	if e != nil {
		log.Fatalf("Error with %s: %v", d, e)
	}

}
