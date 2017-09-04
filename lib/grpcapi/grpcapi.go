package grpcapi

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	app "lekovr/exam/counter"
	"lekovr/exam/counter/setup"
	"lekovr/exam/lib/iface/kvstore"
	"lekovr/exam/lib/iface/logger"
	pb "lekovr/exam/lib/proto/counter"
)

type Config struct {
	Number            int64 `long:"init_number" default:"0" description:"Initial number"`
	Step              int64 `long:"init_step"   default:"1" description:"Increment step"`
	Limit             int64 `long:"init_limit"  default:"100" description:"Increment loop limit"`
	StrictStoreErrors bool  `long:"store_strict"  description:"Do not ignore store errors"`
}

type CounterService struct {
	counter     *app.Counter
	log         logger.Entry
	store       kvstore.Store
	storeStrict bool
}

func NewAPI(log logger.Entry, store kvstore.Store, cfg Config) (*CounterService, error) {

	log.WithField("config", cfg).Debug("Create API")

	sets, err := store.GetSettings()
	if err != nil {
		return nil, err
	} else if sets == nil {
		sets = &setup.Settings{Step: cfg.Step, Limit: cfg.Limit}
		log.WithField("settings", *sets).Debug("Load Settings from opts")
		store.SetSettings(sets)
	} else {
		log.WithField("settings", *sets).Debug("Load Settings from db")
	}

	number, err := store.GetNumber()
	if err != nil {
		return nil, err
	} else if number == nil {
		number = &cfg.Number
		log.WithField("number", *number).Debug("Load Number from opts")
		store.SetNumber(number)
	} else {
		log.WithField("number", *number).Debug("Load Number from db")
	}

	c, err := app.NewCounter(*sets, *number)
	if err != nil {
		return nil, err
	}
	service := CounterService{counter: c, store: store, log: log, storeStrict: cfg.StrictStoreErrors}
	log.Info("API created")
	return &service, nil
}

func (s *CounterService) Close() {
	if s.store != nil {
		s.store.Close()
	}
}

func (s *CounterService) GetNumber(ctx context.Context, in *empty.Empty) (*pb.Number, error) {
	var pbn pb.Number
	number, err := s.counter.GetNumber()
	if err != nil {
		s.log.Errorf("GetNumber error: %+v", err)
		return &pbn, grpc.Errorf(codes.Internal, err.Error())
	}

	pbn.Number = number
	s.log.WithField("number", number).Debug("GetNumber")
	return &pbn, err
}

func (s *CounterService) IncrementNumber(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	number, err := s.counter.IncrementNumber()
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}
	s.log.WithField("number", number).Debug("IncrementNumber")

	err = s.store.SetNumber(&number)
	if err != nil {
		s.log.WithField("number", number).Warnf("Number store error: %+v", err)
		if s.storeStrict {
			return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
		}
	}

	return &empty.Empty{}, nil
}

func (s *CounterService) SetSettings(ctx context.Context, in *pb.Settings) (*empty.Empty, error) {

	sets := setup.Settings{Step: in.Step, Limit: in.Limit}
	s.log.WithField("settings", sets).Debug("SetSettings")

	err := s.counter.SetSettings(sets)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
	err = s.store.SetSettings(&sets)
	if err != nil {
		s.log.WithField("settings", sets).Warnf("Settings store error: %+v", err)
		if s.storeStrict {
			return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
		}
	}
	return &empty.Empty{}, nil
}

func (s *CounterService) GetSettings(ctx context.Context, in *empty.Empty) (*pb.Settings, error) {

	se, err := s.counter.GetSettings()
	if err != nil {
		s.log.Errorf("GetSettings error: %+v", err)
		return &pb.Settings{}, grpc.Errorf(codes.Internal, err.Error())
	}
	s.log.WithField("settings", se).Debug("GetSettings")
	pbs := pb.Settings{Step: se.Step, Limit: se.Limit}
	return &pbs, nil
}
