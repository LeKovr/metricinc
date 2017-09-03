package main

import (
	"net"

	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"
	"github.com/jessevdk/go-flags"

	pb "lekovr/exam/lib/proto/counter"
	"lekovr/exam/lib/boltdb"
	"lekovr/exam/lib/logger"
	"lekovr/exam/lib/server"

	logiface "lekovr/exam/lib/struct/logger"
)

type Config struct {
	Listen string `long:"listen" default:":50051" description:""`

	Logger logger.Config
	Server server.Config
	Store  boltdb.Config
}

func main() {

	var cfg Config

	_, err := flags.Parse(&cfg)
	if err != nil {
		panic("Parse options error") // error message written already
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}

	log.Infof("Counter server v%s", Version)

	lis, err := net.Listen("tcp", cfg.Listen)
	stopOnError(log, err, "listen socket")

	store, err := boltdb.NewStore(log, cfg.Store)
	stopOnError(log, err, "db connect")
	srv, err := server.NewServer(log, store, cfg.Server)
	stopOnError(log, err, "service create")
	defer srv.Close()

	s := grpc.NewServer()
	pb.RegisterCounterServer(s, srv)
	// Register reflection service on gRPC server.
	//	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic("Service init error: " + err.Error())
	}
}

func stopOnError(log logiface.Entry, e error, d string) {
	if e != nil {
		log.Fatalf("Error with %s: %v", d, e)
	}

}
