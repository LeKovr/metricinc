/*
Client for gRPC counter service

Usage

Run from command line:
 client [OPTIONS] <get | inc | set>

Application Options:
  --connect= Addr and port which server listens (default: :50051)

Help Options:
  -h, --help     Show help message

Available commands:
  get  Get current values
  inc  Increment counter
  set  Set both counter settings (step and limit) using defaults if not given (see 'set -h')

[set command options]
  --step=  Increment step (default: 1)
  --limit= Increment loop limit (default: 100)

Program returns all counter data as json.

Examples

Values from clean new service:
 $ ./client get
 {"number":0,"step":1,"limit":100}

Increment counter:
 $ ./client inc
 {"number":1,"step":1,"limit":100}

Set new step:
 $ ./client set --step 2
 {"number":1,"step":2,"limit":100}

Increment by new step
 $ ./client inc
 {"number":3,"step":2,"limit":100}


*/
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"lekovr/exam/lib/client"
	ilogger "lekovr/exam/lib/iface/logger"
	"lekovr/exam/lib/logger"
	pb "lekovr/exam/lib/proto/counter"
)

// GetCommand holds 'get' command defnition
type GetCommand struct{}

// IncCommand holds 'inc' command defnition
type IncCommand struct{}

// SetCommand holds 'set' command defnition
type SetCommand struct {
	Step  int64 `long:"step"   default:"1" description:"Increment step"`
	Limit int64 `long:"limit"  default:"100" description:"Increment loop limit"`
}

// Config holds connect string and command registry
type Config struct {
	Connect string `long:"connect" default:":50051" description:"Addr and port which server listens"`

	Logger logger.Config `group:"Logging Options"`

	Get GetCommand `command:"get" description:"Get current values"`
	Inc IncCommand `command:"inc" description:"Increment counter"`
	Set SetCommand `command:"set" description:"Set both counter settings (step and limit) using defaults if not given (see 'set -h')"`
}

// Response holds data which any command returns
type Response struct {
	Number int64 `json:"number"`
	Step   int64 `json:"step"`
	Limit  int64 `json:"limit"`
}

var cfg Config

// Execute for 'get' command
func (x *GetCommand) Execute(args []string) error {
	number(cfg, false)
	return nil
}

// Execute for 'inc' command
func (x *IncCommand) Execute(args []string) error {
	number(cfg, true)
	return nil
}

// Execute for 'set' command
func (x *SetCommand) Execute(args []string) error {
	settings(cfg, x.Step, x.Limit)
	return nil
}

// main - program body
func main() {

	p := flags.NewParser(&cfg, flags.Default)
	_, err := p.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(1) // help printed
		} else {
			os.Exit(2) // error message written already
		}
	}

	os.Exit(0)
}

// number - connect to gRPC and get current number, if doInc than after increment it
func number(cfg Config, doInc bool) {
	log, c := open(cfg)
	defer c.Close()
	if doInc {
		_, err := c.Service.IncrementNumber(context.Background(), &empty.Empty{})
		if err != nil {
			log.Fatalf("could not inc: %v", err)
		}
	}
	show(log, c)
}

// settings - connect to gRPC and save settings
func settings(cfg Config, step, limit int64) {
	log, c := open(cfg)
	defer c.Close()
	in := pb.Settings{Step: step, Limit: limit}
	_, err := c.Service.SetSettings(context.Background(), &in)
	if err != nil {
		log.Fatalf("could not set settings: %v", err)
	}
	show(log, c)
}

func open(cfg Config) (ilogger.Entry, *client.Count) {
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic("Logger init error: " + err.Error())
	}
	c, err := client.NewClient(cfg.Connect)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return log, c
}

// show current values of number ans settings as json
func show(log ilogger.Entry, c *client.Count) {
	n, err := c.Service.GetNumber(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not get number: %v", err)
	}
	s, err := c.Service.GetSettings(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not get settings: %v", err)
	}
	ret := Response{Number: n.Number, Step: s.Step, Limit: s.Limit}
	payload, _ := json.Marshal(ret)
	fmt.Println(string(payload))
}
