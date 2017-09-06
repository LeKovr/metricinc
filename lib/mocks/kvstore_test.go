// This file contains test templates only
// No tests yet

package mocks

import (
	setup "lekovr/exam/counter/setup"
	reflect "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestNewMockStore(t *testing.T) {
	type args struct {
		ctrl *gomock.Controller
	}
	tests := []struct {
		name string
		args args
		want *MockStore
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := NewMockStore(tt.args.ctrl); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewMockStore() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_EXPECT(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	tests := []struct {
		name   string
		fields fields
		want   *MockStoreMockRecorder
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		if got := m.EXPECT(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStore.EXPECT() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_Close(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		if err := m.Close(); (err != nil) != tt.wantErr {
			t.Errorf("%q. MockStore.Close() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestMockStoreMockRecorder_Close(t *testing.T) {
	type fields struct {
		mock *MockStore
	}
	tests := []struct {
		name   string
		fields fields
		want   *gomock.Call
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mr := &MockStoreMockRecorder{
			mock: tt.fields.mock,
		}
		if got := mr.Close(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStoreMockRecorder.Close() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_GetNumber(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	tests := []struct {
		name    string
		fields  fields
		want    *int64
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		got, err := m.GetNumber()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. MockStore.GetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. MockStore.GetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStoreMockRecorder_GetNumber(t *testing.T) {
	type fields struct {
		mock *MockStore
	}
	tests := []struct {
		name   string
		fields fields
		want   *gomock.Call
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mr := &MockStoreMockRecorder{
			mock: tt.fields.mock,
		}
		if got := mr.GetNumber(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStoreMockRecorder.GetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_GetSettings(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	tests := []struct {
		name    string
		fields  fields
		want    *setup.Settings
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		got, err := m.GetSettings()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. MockStore.GetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStore.GetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStoreMockRecorder_GetSettings(t *testing.T) {
	type fields struct {
		mock *MockStore
	}
	tests := []struct {
		name   string
		fields fields
		want   *gomock.Call
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mr := &MockStoreMockRecorder{
			mock: tt.fields.mock,
		}
		if got := mr.GetSettings(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStoreMockRecorder.GetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_SetNumber(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	type args struct {
		arg0 *int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		if err := m.SetNumber(tt.args.arg0); (err != nil) != tt.wantErr {
			t.Errorf("%q. MockStore.SetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestMockStoreMockRecorder_SetNumber(t *testing.T) {
	type fields struct {
		mock *MockStore
	}
	type args struct {
		arg0 interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *gomock.Call
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mr := &MockStoreMockRecorder{
			mock: tt.fields.mock,
		}
		if got := mr.SetNumber(tt.args.arg0); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStoreMockRecorder.SetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestMockStore_SetSettings(t *testing.T) {
	type fields struct {
		ctrl     *gomock.Controller
		recorder *MockStoreMockRecorder
	}
	type args struct {
		arg0 *setup.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		m := &MockStore{
			ctrl:     tt.fields.ctrl,
			recorder: tt.fields.recorder,
		}
		if err := m.SetSettings(tt.args.arg0); (err != nil) != tt.wantErr {
			t.Errorf("%q. MockStore.SetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestMockStoreMockRecorder_SetSettings(t *testing.T) {
	type fields struct {
		mock *MockStore
	}
	type args struct {
		arg0 interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *gomock.Call
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		mr := &MockStoreMockRecorder{
			mock: tt.fields.mock,
		}
		if got := mr.SetSettings(tt.args.arg0); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MockStoreMockRecorder.SetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
