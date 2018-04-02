package client

import (
	"net"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/LeKovr/metricinc/counter/setup"
	"github.com/LeKovr/metricinc/lib/grpcapi"
	"github.com/LeKovr/metricinc/lib/logger"
	"github.com/LeKovr/metricinc/lib/mocks"
	pb "github.com/LeKovr/metricinc/lib/proto/counter"
)

func TestNewClient(t *testing.T) {
	type args struct {
		address string
	}
	type mock []*gomock.Call

	var zero int64

	tests := []struct {
		name    string
		args    []grpc.DialOption
		want    *pb.Number
		wantErr bool
		mocks   func(st *mocks.MockStore) mock
	}{
		// Test cases.
		{
			name: "NewClient",
			args: []grpc.DialOption{grpc.WithInsecure()},
			want: &pb.Number{Number: 0},
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: 1, Limit: 10}, nil),
					st.EXPECT().GetNumber().Return(&zero, nil),
					st.EXPECT().Close(),
				}
			},
		},
		{
			name:    "NewClient connect error",
			args:    []grpc.DialOption{},
			wantErr: true,
			mocks: func(st *mocks.MockStore) mock {
				return mock{
					st.EXPECT().GetSettings().Return(&setup.Settings{Step: 1, Limit: 10}, nil),
					st.EXPECT().GetNumber().Return(&zero, nil),
					st.EXPECT().Close(),
				}
			},
		},
	}

	cfg := struct {
		Connect string
		Logger  logger.Config
		API     grpcapi.Config
	}{
		Connect: ":50500",
		Logger:  logger.Config{Level: "fatal", UseStdOut: true},
		API:     grpcapi.Config{Number: 0, Step: 1, Limit: 10},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lis, err := net.Listen("tcp", cfg.Connect)
	if err != nil {
		t.Errorf("%q. NewClient() Listen error = %v", "Init", err)
	}

	for _, tt := range tests {

		// Logger
		log, err := logger.NewLogger(cfg.Logger)
		if err != nil {
			t.Errorf("%q. NewClient() Logger error = %v", tt.name, err)
			continue
		}

		// Store
		store := mocks.NewMockStore(ctrl)
		if tt.mocks != nil {
			mo := tt.mocks(store)
			gomock.InOrder(mo...)
		}

		// gRPC API
		api, err := grpcapi.NewAPI(log, store, cfg.API)
		defer api.Close()

		s := grpc.NewServer()
		pb.RegisterCounterServer(s, api)

		// Serve
		go func() { s.Serve(lis) }()
		defer s.Stop()

		c, err := NewClient(cfg.Connect, tt.args...)

		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewClient() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil {
			defer c.Close()
			got, _ := c.Service.GetNumber(context.TODO(), &empty.Empty{})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q. NewClient() = %v, want %v", tt.name, got, tt.want)
			}

		}

	}
}
