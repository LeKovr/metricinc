// Copyright 2017 Alexey Kovrizhkin <lekovr@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"

	"lekovr/exam/lib/boltdb"
	"lekovr/exam/lib/grpcapi"
	"lekovr/exam/lib/logger"
	pb "lekovr/exam/lib/proto/counter"

	logiface "lekovr/exam/lib/iface/logger"
)

// Config holds program options
type Config struct {
	Listen string `long:"listen" default:":50051" description:"Addr and port which server listens at"`

	Logger logger.Config  `group:"Logging Options"`
	API    grpcapi.Config `group:"API Options"`
	Store  boltdb.Config  `group:"Storage Options"`
}

// main - program body
func main() {

	// Parse options
	var cfg Config
	_, err := flags.Parse(&cfg)
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(1) // help printed
		} else {
			os.Exit(2) // error message written already
		}
	}

	// Logger
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}

	log.Infof("Counter server v%s", Version)

	// Listener
	log.WithField("addr", cfg.Listen).Debug("Create listener")
	lis, err := net.Listen("tcp", cfg.Listen)
	stopOnError(log, err, "Listen socket")

	// Store
	store, err := boltdb.NewStore(log, cfg.Store)
	stopOnError(log, err, "DB connect")

	// gRPC API
	api, err := grpcapi.NewAPI(log, store, cfg.API)
	stopOnError(log, err, "Service create")
	defer api.Close()

	s := grpc.NewServer()
	pb.RegisterCounterServer(s, api)

	// Gracefull shutdown
	// http://stackoverflow.com/questions/18106749/golang-catch-signals
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Infof("Got signal: %v", sig)
		api.Close()
		s.GracefulStop()
		os.Exit(0)
	}()

	// Serve
	if err := s.Serve(lis); err != nil {
		panic("Server init error: " + err.Error())
	}
}

// stopOnError used internally for fatal errors checking
func stopOnError(log logiface.Entry, err error, info string) {
	if err != nil {
		log.Fatalf("Error with %s: %v", info, err)
	}

}
