package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	app "lekovr/exam/counter"
	"lekovr/exam/counter/setup"
	pb "lekovr/exam/lib/proto/counter"
	"lekovr/exam/lib/struct/kvstore"
	"lekovr/exam/lib/struct/logger"
)

type Config struct {
	//	NumberKey []byte
	Number int64 `long:"init_number" default:"0" description:"Initial number"`
	Step   int64 `long:"init_step"   default:"1" description:"Increment step"`
	Limit  int64 `long:"init_limit"  default:"100" description:"Increment loop limit"`
}

type CounterService struct {
	counter *app.Counter
	store   kvstore.Store
	log     logger.Entry
}

func NewServer(log logger.Entry, store kvstore.Store, cfg Config) (*CounterService, error) {

	sets, err := store.GetSettings()
	if err != nil {
		return nil, err
	} else if sets == nil {
		sets = &setup.Settings{Step: cfg.Step, Limit: cfg.Limit}
		log.Debugf("Load Settings from opts: %+v", sets)
		store.SetSettings(sets)
	}

	number, err := store.GetNumber()
	if err != nil {
		return nil, err
	} else if number == nil {
		number = &cfg.Number
		store.SetNumber(number)
	}

	c, err := app.NewCounter(*sets, *number)
	if err != nil {
		return nil, err
	}
	service := CounterService{counter: c, store: store, log: log}
	log.Infof("Got service: %+v", service)
	return &service, nil
}

func (s *CounterService) Close() {
	if s.store != nil {
		s.store.Close()
	}
}

func (s *CounterService) GetNumber(ctx context.Context, in *empty.Empty) (*pb.Number, error) {
	var i pb.Number
	n, err := s.counter.GetNumber()
	i.Number = n
	return &i, err
}

func (s *CounterService) IncrementNumber(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	val, err := s.counter.IncrementNumber()
	if err != nil {
		return &empty.Empty{}, err
	}
	s.log.Debugf("Ready to SetNumber: %d", val)

	err = s.store.SetNumber(&val)
	if err != nil {
		s.log.Debugf("ER01: %+v", err)
	}

	return &empty.Empty{}, err
}

func (s *CounterService) SetSettings(ctx context.Context, in *pb.SettingsRequest) (*pb.SettingsReply, error) {

	se := setup.Settings{Step: in.Step, Limit: in.Limit}

	var ret *pb.SettingsReply
	err := s.counter.SetSettings(se)
	if err != nil {
		ret = &pb.SettingsReply{Code: 1, Message: err.Error()}
	} else {
		err = s.store.SetSettings(&se)
		if err == nil {
			ret = &pb.SettingsReply{Code: 0, Message: "Ok"}
		} else {
			ret = &pb.SettingsReply{Code: 1, Message: err.Error()}
		}
	}

	return ret, err
}
