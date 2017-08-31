package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	pb "lekovr/exam/counter"
)

type CounterService struct {
	Value int64
}

func NewServer() (*CounterService, error) {
	service := CounterService{Value: 0}
	return &service, nil
}

func (s *CounterService) GetNumber(ctx context.Context, in *empty.Empty) (*pb.Value, error) {
	var i pb.Value
	i.Value = s.Value
	return &i, nil
}

func (s *CounterService) IncrementNumber(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	//  var i pb.Value;
	s.Value++
	return &empty.Empty{}, nil
}

func (s *CounterService) SetSettings(ctx context.Context, in *pb.SettingsRequest) (*pb.SettingsReply, error) {
	ret := pb.SettingsReply{Code: 1, Message: "Ok"}
	return &ret, nil
}
