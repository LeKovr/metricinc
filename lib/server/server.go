package server

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	pb "lekovr/exam/counter"
	"lekovr/exam/lib/struct/kvstore"
	"lekovr/exam/lib/struct/server"
)

type Config struct {
	//	NumberKey []byte
	Number int64 `long:"init_number" default:"0" description:""`
	Step   int64 `long:"init_step"   default:"1" description:""`
	Limit  int64 `long:"init_limit"  default:"100" description:""`
}

type CounterService struct {
	Number   int64
	Settings server.Settings
	store    kvstore.Store
}

func NewServer(cfg Config, store kvstore.Store) (*CounterService, error) {

	sets, err := store.GetSettings()
	if err != nil {
		return nil, err
	} else if sets == nil {
		sets = &server.Settings{Step: cfg.Step, Limit: cfg.Limit}
		store.SetSettings(sets)
	}
	number, err := store.GetNumber()
	if err != nil {
		return nil, err
	} else if number == nil {
		number = &cfg.Number
		store.SetNumber(number)
	}

	service := CounterService{Number: *number, Settings: *sets, store: store}
	log.Println("Got service")
	return &service, nil
}

func (s *CounterService) Close() {
	if s.store != nil {
		s.store.Close()
	}
}

func (s *CounterService) GetNumber(ctx context.Context, in *empty.Empty) (*pb.Number, error) {
	var i pb.Number
	i.Number = s.Number
	return &i, nil
}

func (s *CounterService) IncrementNumber(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	//  var i pb.Number;
	if s.Number > (s.Settings.Limit - s.Settings.Step) {
		s.Number = (s.Number + s.Settings.Step - s.Settings.Limit) // TODO: Number+Step может быть переполнение
	} else {
		s.Number += s.Settings.Step
	}
	log.Printf("Ready to SetNumber: %d", s.Number)

	err := s.store.SetNumber(&s.Number)
	if err != nil {
		log.Printf("ER01: %+v", err)
	}

	return &empty.Empty{}, err
}

func (s *CounterService) SetSettings(ctx context.Context, in *pb.SettingsRequest) (*pb.SettingsReply, error) {
	var ret *pb.SettingsReply
	s.Settings = server.Settings{Step: in.Step, Limit: in.Limit}
	err := s.store.SetSettings(&s.Settings)
	if err == nil {
		ret = &pb.SettingsReply{Code: 0, Message: "Ok"}
	} else {
		ret = &pb.SettingsReply{Code: 1, Message: err.Error()}
	}
	return ret, err
}
