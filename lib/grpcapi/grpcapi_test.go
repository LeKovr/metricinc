package grpcapi

import (
	app "lekovr/exam/counter"
	"lekovr/exam/counter/setup"
	//	"lekovr/exam/lib/iface/kvstore"
	"lekovr/exam/lib/logger"
	//	pb "lekovr/exam/lib/proto/counter"
	"reflect"
	"testing"

	//"github.com/Sirupsen/logrus"
	//	"github.com/Sirupsen/logrus/hooks/test"

	gomock "github.com/golang/mock/gomock"
	"lekovr/exam/lib/mock_kvstore"
	//	"github.com/golang/protobuf/ptypes/empty"
	//	"golang.org/x/net/context"
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
	tests := []struct {
		name    string
		args    args
		want    *CounterService
		wantErr bool
	}{
		// Test cases.
		{
			name: "Simple new",
			args: args{cfg: Config{Number: 0, Step: 1, Limit: 10}},
			//			want: &CounterService{},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var cfg struct {
		Logger logger.Config `group:"Logging Options"`
	}
	cfg.Logger.Level = "warn"
	cfg.Logger.UseStdOut = true
	log, _ := logger.NewLogger(cfg.Logger)

	for _, tt := range tests {
		//log, hook := test.NewNullLogger()
		store := mock_kvstore.NewMockStore(ctrl)
		store.EXPECT().GetSettings()
		sets := &setup.Settings{Step: tt.args.cfg.Step, Limit: tt.args.cfg.Limit}
		store.EXPECT().SetSettings(sets)
		store.EXPECT().GetNumber()
		store.EXPECT().SetNumber(&tt.args.cfg.Number)
		got, err := NewAPI(log, store, tt.args.cfg)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewAPI() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		counter, _ := app.NewCounter(*sets, tt.args.cfg.Number)

		tt.want = &CounterService{counter: counter, log: log, store: store}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewAPI() = %v, want %v", tt.name, got, tt.want)
		}

		//		hook.Reset()
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

	store := mock_kvstore.NewMockStore(ctrl)
	sets := &setup.Settings{Step: 1, Limit: 10}

	counter, _ := app.NewCounter(*sets, 0)

	s := &CounterService{
		counter: counter,
		store:   store,
		log:     log,
	}
	store.EXPECT().Close()
	s.Close()
}

/*
func TestCounterService_GetNumber(t *testing.T) {
	type fields struct {
		counter     *app.Counter
		log         logger.Entry
		store       kvstore.Store
		storeStrict bool
	}
	type args struct {
		ctx context.Context
		in  *empty.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Number
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		s := &CounterService{
			counter:     tt.fields.counter,
			log:         tt.fields.log,
			store:       tt.fields.store,
			storeStrict: tt.fields.storeStrict,
		}
		got, err := s.GetNumber(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.GetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.GetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounterService_IncrementNumber(t *testing.T) {
	type fields struct {
		counter     *app.Counter
		log         logger.Entry
		store       kvstore.Store
		storeStrict bool
	}
	type args struct {
		ctx context.Context
		in  *empty.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		s := &CounterService{
			counter:     tt.fields.counter,
			log:         tt.fields.log,
			store:       tt.fields.store,
			storeStrict: tt.fields.storeStrict,
		}
		got, err := s.IncrementNumber(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.IncrementNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.IncrementNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounterService_SetSettings(t *testing.T) {
	type fields struct {
		counter     *app.Counter
		log         logger.Entry
		store       kvstore.Store
		storeStrict bool
	}
	type args struct {
		ctx context.Context
		in  *pb.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		s := &CounterService{
			counter:     tt.fields.counter,
			log:         tt.fields.log,
			store:       tt.fields.store,
			storeStrict: tt.fields.storeStrict,
		}
		got, err := s.SetSettings(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.SetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.SetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounterService_GetSettings(t *testing.T) {
	type fields struct {
		counter     *app.Counter
		log         logger.Entry
		store       kvstore.Store
		storeStrict bool
	}
	type args struct {
		ctx context.Context
		in  *empty.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Settings
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		s := &CounterService{
			counter:     tt.fields.counter,
			log:         tt.fields.log,
			store:       tt.fields.store,
			storeStrict: tt.fields.storeStrict,
		}
		got, err := s.GetSettings(tt.args.ctx, tt.args.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. CounterService.GetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. CounterService.GetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
*/
