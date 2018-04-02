package grpcapi

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	app "github.com/LeKovr/metricinc/counter"
	"github.com/LeKovr/metricinc/counter/setup"
	"github.com/LeKovr/metricinc/lib/iface/kvstore"
	"github.com/LeKovr/metricinc/lib/iface/logger"
	pb "github.com/LeKovr/metricinc/lib/proto/counter"
)

// Config is a program flags group used in constructor
type Config struct {
	Number            int64 `long:"init_number" default:"0" description:"Initial number"`
	Step              int64 `long:"init_step"   default:"1" description:"Increment step"`
	Limit             int64 `long:"init_limit"  default:"100" description:"Increment loop limit"`
	StrictStoreErrors bool  `long:"store_strict"  description:"Do not ignore store errors"`
}

// CounterService holds object internals
type CounterService struct {
	counter     *app.Counter
	log         logger.Entry
	store       kvstore.Store
	storeStrict bool
}

// NewAPI creates an API object
func NewAPI(log logger.Entry, store kvstore.Store, cfg Config) (*CounterService, error) {

	log.WithField("config", cfg).Debug("Create API")

	newSettings := false
	newNumber := false

	// Get settings from db or use defaults from cfg
	sets, err := store.GetSettings()
	if err != nil {
		return nil, err
	} else if sets == nil {
		sets = &setup.Settings{Step: cfg.Step, Limit: cfg.Limit}
		newSettings = true
		log.WithField("settings", *sets).Debug("Load Settings from opts")
	} else {
		log.WithField("settings", *sets).Debug("Load Settings from db")
	}

	// Get number from db or use defaults from cfg
	number, err := store.GetNumber()
	if err != nil {
		return nil, err
	} else if number == nil {
		number = &cfg.Number
		newNumber = true
		log.WithField("number", *number).Debug("Load Number from opts")
	} else {
		log.WithField("number", *number).Debug("Load Number from db")
	}

	// Create counter
	c, err := app.NewCounter(sets, *number)
	if err != nil {
		return nil, err
	}
	// All data ready, store if new
	if newSettings {
		store.SetSettings(sets)
	}
	if newNumber {
		store.SetNumber(number)
	}
	service := CounterService{counter: c, store: store, log: log, storeStrict: cfg.StrictStoreErrors}
	log.Info("API created")
	return &service, nil
}

// GetSettings reads settings from counter
func (s *CounterService) GetSettings(ctx context.Context, in *empty.Empty) (*pb.Settings, error) {

	sets, err := s.counter.GetSettings()
	if err != nil {
		s.log.Errorf("GetSettings error: %+v", err)
		return &pb.Settings{}, grpc.Errorf(codes.Internal, err.Error())
	}
	s.log.WithField("settings", sets).Debug("GetSettings")
	pbs := pb.Settings{Step: sets.Step, Limit: sets.Limit}
	return &pbs, nil
}

// SetSettings sets settings to counter and stores them in database
func (s *CounterService) SetSettings(ctx context.Context, in *pb.Settings) (*empty.Empty, error) {

	sets := setup.Settings{Step: in.Step, Limit: in.Limit}
	s.log.WithField("settings", sets).Debug("SetSettings")

	// Set in counter
	err := s.counter.SetSettings(&sets)
	if err != nil {
		s.log.WithField("settings", sets).Warnf("Settings set error: %+v", err)
		return &empty.Empty{}, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	// Set in store
	err = s.store.SetSettings(&sets)
	if err != nil {
		s.log.WithField("settings", sets).Warnf("Settings store error: %+v", err)
		if s.storeStrict {
			return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
		}
	}
	return &empty.Empty{}, nil
}

// GetNumber reads number from counter
func (s *CounterService) GetNumber(ctx context.Context, in *empty.Empty) (*pb.Number, error) {
	var pbn pb.Number
	number, err := s.counter.GetNumber()
	if err != nil {
		s.log.Errorf("GetNumber error: %+v", err)
		return &pbn, grpc.Errorf(codes.Internal, err.Error())
	}

	pbn.Number = *number
	s.log.WithField("number", *number).Debug("GetNumber")
	return &pbn, err
}

// IncrementNumber increments counter and stores new number in database
func (s *CounterService) IncrementNumber(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {

	// Set in counter
	number, err := s.counter.IncrementNumber()
	if err != nil {
		s.log.WithField("number", number).Warnf("Number inc error: %+v", err)
		return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
	}

	s.log.WithField("number", number).Debug("IncrementNumber")

	// Set in store
	err = s.store.SetNumber(number)
	if err != nil {
		s.log.WithField("number", number).Warnf("Number store error: %+v", err)
		if s.storeStrict {
			return &empty.Empty{}, grpc.Errorf(codes.Internal, err.Error())
		}
	}
	return &empty.Empty{}, nil
}

// Close logs final values as warning and closes store.
func (s *CounterService) Close() {
	if s.store != nil {
		number, _ := s.counter.GetNumber()
		sets, _ := s.counter.GetSettings()
		s.log.WithField("settings", *sets).WithField("number", *number).Warn("Final state")
		s.store.Close()
	}
}
