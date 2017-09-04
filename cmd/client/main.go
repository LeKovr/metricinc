package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"lekovr/exam/lib/client"
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
	Connect string     `long:"connect" default:":50051" description:"Addr and port which server listens"`
	Get     GetCommand `command:"get" description:"Get current values"`
	Inc     IncCommand `command:"inc" description:"Increment counter"`
	Set     SetCommand `command:"set" description:"Set both counter settings (step and limit) using defaults if not given (see 'set -h')"`
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
	number(cfg.Connect, false)
	return nil
}

// Execute for 'inc' command
func (x *IncCommand) Execute(args []string) error {
	number(cfg.Connect, true)
	return nil
}

// Execute for 'set' command
func (x *SetCommand) Execute(args []string) error {
	settings(cfg.Connect, x.Step, x.Limit)
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
func number(connect string, doInc bool) {
	c, err := client.NewServer(connect)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer c.Close()
	if doInc {
		_, err = c.Service.IncrementNumber(context.Background(), &empty.Empty{})
		if err != nil {
			log.Fatalf("could not inc: %v", err)
		}
	}
	show(c)
}

// settings - connect to gRPC and save settings
func settings(connect string, step, limit int64) {
	c, err := client.NewServer(connect)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer c.Close()
	in := pb.Settings{Step: step, Limit: limit}
	_, err = c.Service.SetSettings(context.Background(), &in)
	if err != nil {
		log.Fatalf("could not set settings: %v", err)
	}
	show(c)
}

// show current values of number ans settings as json
func show(c *client.Count) {
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
