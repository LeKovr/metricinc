package grpcapi

import (
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	app "lekovr/exam/counter"
	"lekovr/exam/counter/setup"
	"lekovr/exam/lib/logger"
	"lekovr/exam/lib/mocks"
	pb "lekovr/exam/lib/proto/counter"
)

/*
import(
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestSomething(t*testing.T){
  logger, hook := test.NewNullLogger()
  logger.Error("Helloerror")

  assert.Equal(t, 1, len(hook.Entries))
  assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
  assert.Equal(t, "Helloerror", hook.LastEntry().Message)

  hook.Reset()
  assert.Nil(t, hook.LastEntry())
}
*/

func TestNewAPI(t *testing.T) {
	type args struct {
		cfg Config
	}
	type mock []*gomock.Call

	var zero int64
	var five int64 = 5

	tests := []struct {
		name    string
		args    args
		want    *CounterService
		wantErr bool
		sets    *setup.Settings
		mocks   func(st *mocks.MockStore) mock
	}{
		// Test cases.
		{
			name: "Empty store",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			// want: - calculated inside loop
			sets: &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings(),
					st.EXPECT().GetNumber(),
					st.EXPECT().SetSettings(&setup.Settings{Step: 1, Limit: 10}),
					st.EXPECT().SetNumber(&zero),
				}
			},
		},

		{
			name: "Bad config",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			// want: - calculated inside loop
			wantErr: true,
			sets:    &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: -1, Limit: 10}, nil),
					st.EXPECT().GetNumber().Return(&zero, nil),
				}
			},
		},
		{
			name: "Store with settings",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			// want: - calculated inside loop
			sets: &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: 1, Limit: 10}, nil),
					st.EXPECT().GetNumber(),
					st.EXPECT().SetNumber(&zero),
				}
			},
		},
		{
			name:    "Store with settings error",
			args:    args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			wantErr: true,
			sets:    &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(nil, errors.New("read error")),
				}
			},
		},
		{
			name: "Store with settings & number",
			args: args{cfg: Config{Number: five, Step: 1, Limit: 10}},
			// want: - calculated inside loop
			sets: &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: 1, Limit: 10}, nil),
					st.EXPECT().GetNumber().Return(&five, nil),
				}
			},
		},
		{
			name:    "Store with number error",
			args:    args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			wantErr: true,
			sets:    &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: 1, Limit: 10}, nil),
					st.EXPECT().GetNumber().Return(nil, errors.New("read error")),
				}
			},
		},
	}

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "warn"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		store := mocks.NewMockStore(ctrl)
		if tt.mocks != nil {
			mo := tt.mocks(store)
			gomock.InOrder(mo...)
		}

		got, err := NewAPI(log, store, tt.args.cfg)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewAPI() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		counter, _ := app.NewCounter(tt.sets, tt.args.cfg.Number)
		if !tt.wantErr {
			tt.want = &CounterService{counter: counter, log: log, store: store}
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewAPI() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounterService_Close(t *testing.T) {

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "error"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	sets := &setup.Settings{Step: 1, Limit: 10}

	counter, _ := app.NewCounter(sets, 0)

	s := &CounterService{
		counter: counter,
		store:   store,
		log:     log,
	}
	store.EXPECT().Close()
	s.Close()
}

func TestCounterService_GetNumber(t *testing.T) {

	type args struct {
		cfg Config
		in  *empty.Empty
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.Number
		wantErr bool
		sets    *setup.Settings
	}{
		// Test cases.
		{
			name: "GetNumber ok",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			want: &pb.Number{Number: 0},
			sets: &setup.Settings{Step: 1, Limit: 10},
		},
		{
			name:    "GetNumber error",
			wantErr: true,
		},
	}

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "fatal"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {

		store := mocks.NewMockStore(ctrl)
		var counter *app.Counter
		if tt.wantErr {
			counter = &app.Counter{}
		} else {
			counter, _ = app.NewCounter(tt.sets, tt.args.cfg.Number)
		}
		s := &CounterService{
			counter: counter,
			log:     log,
			store:   store,
		}
		got, err := s.GetNumber(context.TODO(), tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.GetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.GetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounterService_IncrementNumber(t *testing.T) {

	type args struct {
		cfg Config
		in  *empty.Empty
	}

	type mock []*gomock.Call

	var five int64 = 5

	tests := []struct {
		name         string
		args         args
		want         *pb.Number
		wantErr      bool
		wantStoreErr bool
		sets         *setup.Settings
		mocks        func(st *mocks.MockStore) mock
		storeStrict  bool
	}{
		// Test cases.
		{
			name: "IncrementNumber ok",
			args: args{cfg: Config{Number: 4}},
			sets: &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetNumber(&five),
				}
			},
		},
		{
			name:    "IncrementNumber counter error",
			wantErr: true,
		},
		{
			name: "IncrementNumber store error not fired",
			args: args{cfg: Config{Number: 4}},
			sets: &setup.Settings{Step: 1, Limit: 10},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetNumber(&five).Return(errors.New("write error")),
				}
			},
		},
		{
			name:         "IncrementNumber store strict error",
			args:         args{cfg: Config{Number: 4}},
			sets:         &setup.Settings{Step: 1, Limit: 10},
			wantStoreErr: true,
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetNumber(&five).Return(errors.New("write error")),
				}
			},
			storeStrict: true,
		},
	}

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "fatal"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {

		store := mocks.NewMockStore(ctrl)
		var counter *app.Counter
		if tt.wantErr {
			counter = &app.Counter{}
		} else {
			counter, _ = app.NewCounter(tt.sets, tt.args.cfg.Number)
		}
		if tt.mocks != nil {
			mo := tt.mocks(store)
			gomock.InOrder(mo...)
		}

		s := &CounterService{
			counter:     counter,
			log:         log,
			store:       store,
			storeStrict: tt.storeStrict,
		}
		_, err := s.IncrementNumber(context.TODO(), tt.args.in)
		if (err != nil) != (tt.wantErr || tt.wantStoreErr) {
			t.Errorf("%q. CounterService.IncrementNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr || tt.wantStoreErr)
			continue
		}
	}
}

func TestCounterService_SetSettings(t *testing.T) {

	type args struct {
		cfg Config
		in  *pb.Settings
	}

	type mock []*gomock.Call

	in := pb.Settings{Step: 2, Limit: 20}
	sets := setup.Settings{Step: 2, Limit: 20}

	tests := []struct {
		name         string
		args         args
		want         *pb.Number
		wantErr      bool
		wantStoreErr bool
		sets         *setup.Settings
		mocks        func(st *mocks.MockStore) mock
		storeStrict  bool
	}{
		// Test cases.
		{
			name: "SetSettings ok",
			args: args{cfg: Config{Number: 0}, in: &in},
			sets: &sets,
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetSettings(&sets),
				}
			},
		},
		{
			name:    "SetSettings counter error",
			args:    args{in: &pb.Settings{Step: 200, Limit: 20}},
			sets:    &sets,
			wantErr: true,
		},
		{
			name: "SetSettings store error not fired",
			args: args{cfg: Config{Number: 0}, in: &in},
			sets: &sets,
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetSettings(&sets).Return(errors.New("write error")),
				}
			},
		},
		{
			name:         "SetSettings store strict error",
			args:         args{cfg: Config{Number: 0}, in: &in},
			sets:         &sets,
			wantStoreErr: true,
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().SetSettings(&sets).Return(errors.New("write error")),
				}
			},
			storeStrict: true,
		},
	}

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "fatal"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {

		store := mocks.NewMockStore(ctrl)
		counter, _ := app.NewCounter(tt.sets, tt.args.cfg.Number)
		if tt.mocks != nil {
			mo := tt.mocks(store)
			gomock.InOrder(mo...)
		}

		s := &CounterService{
			counter:     counter,
			log:         log,
			store:       store,
			storeStrict: tt.storeStrict,
		}
		_, err := s.SetSettings(context.TODO(), tt.args.in)
		if (err != nil) != (tt.wantErr || tt.wantStoreErr) {
			t.Errorf("%q. CounterService.SetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr || tt.wantStoreErr)
			continue
		}
	}
}

func TestCounterService_GetSettings(t *testing.T) {

	type args struct {
		cfg Config
		in  *empty.Empty
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.Settings
		wantErr bool
		sets    *setup.Settings
	}{
		// Test cases.
		{
			name: "GetSettings ok",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			want: &pb.Settings{Step: 1, Limit: 10},
			sets: &setup.Settings{Step: 1, Limit: 10},
		},
		{
			name:    "GetSettings error",
			wantErr: true,
		},
	}

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "fatal"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {

		store := mocks.NewMockStore(ctrl)
		var counter *app.Counter
		if tt.wantErr {
			counter = &app.Counter{}
		} else {
			counter, _ = app.NewCounter(tt.sets, tt.args.cfg.Number)
		}
		s := &CounterService{
			counter: counter,
			log:     log,
			store:   store,
		}
		got, err := s.GetSettings(context.TODO(), tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.GetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.Getsettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
