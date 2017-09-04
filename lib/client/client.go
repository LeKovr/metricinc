package client

import (
	"google.golang.org/grpc"

	pb "lekovr/exam/lib/proto/counter"
)

type Count struct {
	Service pb.CounterClient
	conn    *grpc.ClientConn
}

func NewServer(address string) (*Count, error) {

	c := Count{}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.Service = pb.NewCounterClient(conn)

	return &c, nil
}

func (c *Count) Close() {

	c.conn.Close()
}
