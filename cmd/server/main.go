package main

import (
	"net"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"

	"lekovr/exam/lib/boltdb"
	"lekovr/exam/lib/grpcapi"
	"lekovr/exam/lib/logger"
	pb "lekovr/exam/lib/proto/counter"

	logiface "lekovr/exam/lib/iface/logger"
)

type Config struct {
	Listen string `long:"listen" default:":50051" description:"Addr and port which server listens"`

	Logger logger.Config  `group:"Logging Options"`
	API    grpcapi.Config `group:"API Options"`
	Store  boltdb.Config  `group:"Storage Options"`
}

func main() {

	var cfg Config

	_, err := flags.Parse(&cfg)
	if err != nil {
		panic("Program aborted" + err.Error()) // error message written already
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}

	log.Infof("Counter server v%s", Version)

	log.WithField("addr", cfg.Listen).Debug("Create listener")
	lis, err := net.Listen("tcp", cfg.Listen)
	stopOnError(log, err, "Listen socket")

	store, err := boltdb.NewStore(log, cfg.Store)
	stopOnError(log, err, "DB connect")

	api, err := grpcapi.NewAPI(log, store, cfg.API)
	stopOnError(log, err, "Service create")
	defer api.Close()

	s := grpc.NewServer()
	pb.RegisterCounterServer(s, api)
	if err := s.Serve(lis); err != nil {
		panic("Server init error: " + err.Error())
	}
}

func stopOnError(log logiface.Entry, e error, d string) {
	if e != nil {
		log.Fatalf("Error with %s: %v", d, e)
	}

}
