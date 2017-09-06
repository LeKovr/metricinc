// This file contains test templates only
// TODO: Write tests

package boltdb

import (
	"lekovr/exam/counter/setup"
	"lekovr/exam/lib/iface/logger"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
)

func TestNewStore(t *testing.T) {
	type args struct {
		log logger.Entry
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := NewStore(tt.args.log, tt.args.cfg)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewStore() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewStore() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestStore_GetSettings(t *testing.T) {
	type fields struct {
		Bucket      []byte
		NumberKey   []byte
		SettingsKey []byte
		db          *bolt.DB
		log         logger.Entry
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
		s := &Store{
			Bucket:      tt.fields.Bucket,
			NumberKey:   tt.fields.NumberKey,
			SettingsKey: tt.fields.SettingsKey,
			db:          tt.fields.db,
			log:         tt.fields.log,
		}
		got, err := s.GetSettings()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Store.GetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Store.GetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestStore_SetSettings(t *testing.T) {
	type fields struct {
		Bucket      []byte
		NumberKey   []byte
		SettingsKey []byte
		db          *bolt.DB
		log         logger.Entry
	}
	type args struct {
		sets *setup.Settings
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
		s := &Store{
			Bucket:      tt.fields.Bucket,
			NumberKey:   tt.fields.NumberKey,
			SettingsKey: tt.fields.SettingsKey,
			db:          tt.fields.db,
			log:         tt.fields.log,
		}
		if err := s.SetSettings(tt.args.sets); (err != nil) != tt.wantErr {
			t.Errorf("%q. Store.SetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestStore_GetNumber(t *testing.T) {
	type fields struct {
		Bucket      []byte
		NumberKey   []byte
		SettingsKey []byte
		db          *bolt.DB
		log         logger.Entry
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
		s := &Store{
			Bucket:      tt.fields.Bucket,
			NumberKey:   tt.fields.NumberKey,
			SettingsKey: tt.fields.SettingsKey,
			db:          tt.fields.db,
			log:         tt.fields.log,
		}
		got, err := s.GetNumber()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Store.GetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Store.GetNumber() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestStore_SetNumber(t *testing.T) {
	type fields struct {
		Bucket      []byte
		NumberKey   []byte
		SettingsKey []byte
		db          *bolt.DB
		log         logger.Entry
	}
	type args struct {
		number *int64
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
		s := &Store{
			Bucket:      tt.fields.Bucket,
			NumberKey:   tt.fields.NumberKey,
			SettingsKey: tt.fields.SettingsKey,
			db:          tt.fields.db,
			log:         tt.fields.log,
		}
		if err := s.SetNumber(tt.args.number); (err != nil) != tt.wantErr {
			t.Errorf("%q. Store.SetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestStore_Close(t *testing.T) {
	type fields struct {
		Bucket      []byte
		NumberKey   []byte
		SettingsKey []byte
		db          *bolt.DB
		log         logger.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		s := &Store{
			Bucket:      tt.fields.Bucket,
			NumberKey:   tt.fields.NumberKey,
			SettingsKey: tt.fields.SettingsKey,
			db:          tt.fields.db,
			log:         tt.fields.log,
		}
		if err := s.Close(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Store.Close() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestMarshalSettings(t *testing.T) {
	type args struct {
		sr *setup.Settings
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := MarshalSettings(tt.args.sr)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. MarshalSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. MarshalSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUnmarshalSettings(t *testing.T) {
	type args struct {
		data []byte
		d    *setup.Settings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := UnmarshalSettings(tt.args.data, tt.args.d); (err != nil) != tt.wantErr {
			t.Errorf("%q. UnmarshalSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
