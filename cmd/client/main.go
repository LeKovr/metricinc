package main

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"lekovr/exam/lib/client"
)

const (
	address = "localhost:50051"
)

func main() {

	c, err := client.NewServer(address)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	r, err := c.Service.GetNumber(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not get number1: %v", err)
	}
	log.Printf("Start counter: %d", r.Number)

	_, err = c.Service.IncrementNumber(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not inc: %v", err)
	}
	_, err = c.Service.IncrementNumber(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not int2: %v", err)
	}
	r, err = c.Service.GetNumber(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not get2: %v", err)
	}
	log.Printf("End counter: %d", r.Number)
}
