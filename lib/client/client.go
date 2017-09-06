/*
Package client is a simple connector to counter gRPC service


*/
package client

import (
	"google.golang.org/grpc"

	pb "lekovr/exam/lib/proto/counter"
)

// Count holds gRPC connection & service
type Count struct {
	Service pb.CounterClient
	conn    *grpc.ClientConn
}

// NewClient creates a gRPC connection & service
func NewClient(address string) (*Count, error) {

	c := Count{}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.conn = conn
	// Attach service
	c.Service = pb.NewCounterClient(conn)

	return &c, nil
}

// Close client connection
func (c *Count) Close() {
	c.conn.Close()
}
