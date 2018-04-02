package counter

import (
	"github.com/LeKovr/metricinc/counter/setup"
	"reflect"
	"testing"
)

func TestNewCounter(t *testing.T) {
	type args struct {
		s      *setup.Settings
		number int64
	}
	tests := []struct {
		name    string
		args    args
		want    *Counter
		wantErr bool
	}{
		// Test cases.
		{
			name:    "start from 0",
			args:    args{s: &setup.Settings{Step: 1, Limit: 10}, number: 0},
			want:    &Counter{settings: &setup.Settings{Step: 1, Limit: 10}, number: 0},
			wantErr: false,
		},
		{
			name:    "start from 1",
			args:    args{s: &setup.Settings{Step: 2, Limit: 20}, number: 1},
			want:    &Counter{settings: &setup.Settings{Step: 2, Limit: 20}, number: 1},
			wantErr: false,
		},
		{
			name:    "error: Step must be less than limit",
			args:    args{s: &setup.Settings{Step: 20, Limit: 2}, number: 1},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error: Step must be positive",
			args:    args{s: &setup.Settings{Step: -20, Limit: 2}, number: 1},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		got, err := NewCounter(tt.args.s, tt.args.number)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewCounter() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewCounter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCounter_GetNumber(t *testing.T) {
	type fields struct {
		number   int64
		settings *setup.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// Test cases.
		{
			name:    "start from 1",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 10}, number: 1},
			want:    1,
			wantErr: false,
		},
		{
			name:    "start from 2",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 10}, number: 2},
			want:    2,
			wantErr: false,
		},
		{
			name:    "error: Constructor required",
			fields:  fields{settings: nil, number: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		c := &Counter{
			number:   tt.fields.number,
			settings: tt.fields.settings,
		}
		got, err := c.GetNumber()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Counter.GetNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil && *got != tt.want {
			t.Errorf("%q. Counter.GetNumber() = %v, want %v", tt.name, got, &tt.want)
		}
	}
}

func TestCounter_IncrementNumber(t *testing.T) {
	type fields struct {
		number   int64
		settings *setup.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// Test cases.
		{
			name:    "new number less than limit",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 3}, number: 1},
			want:    2,
			wantErr: false,
		},
		{
			name:    "new number eq limit-1",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 3}, number: 1},
			want:    2,
			wantErr: false,
		},
		{
			name:    "new number eq limit",
			fields:  fields{settings: &setup.Settings{Step: 2, Limit: 3}, number: 1},
			want:    0,
			wantErr: false,
		},
		{
			name:    "new greater than limit",
			fields:  fields{settings: &setup.Settings{Step: 3, Limit: 5}, number: 4},
			want:    2,
			wantErr: false,
		},
		{
			name:    "error: constructor required",
			fields:  fields{settings: nil, number: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		c := &Counter{
			number:   tt.fields.number,
			settings: tt.fields.settings,
		}
		got, err := c.IncrementNumber()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Counter.IncrementNumber() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if err == nil && *got != tt.want {
			t.Errorf("%q. Counter.IncrementNumber() = %v, want %v", tt.name, got, &tt.want)
		}
	}
}

func TestCounter_SetSettings(t *testing.T) {
	type fields struct {
		number   int64
		settings *setup.Settings
	}
	type args struct {
		s setup.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// Test cases.
		{
			name:    "correct settings",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 3}, number: 1},
			args:    args{s: setup.Settings{Step: 2, Limit: 20}},
			wantErr: false,
		},
		{
			name:    "error: Incorrect settings",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 3}, number: 1},
			args:    args{s: setup.Settings{Step: 20, Limit: 2}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		c := &Counter{
			number:   tt.fields.number,
			settings: tt.fields.settings,
		}
		if err := c.SetSettings(&tt.args.s); (err != nil) != tt.wantErr {
			t.Errorf("%q. Counter.SetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if !tt.wantErr && !reflect.DeepEqual(c.settings, &tt.args.s) {
			t.Errorf("%q. Counter.SetSettings() = %v, want %v", tt.name, c.settings, &tt.args.s)
		}

	}
}

func TestCounter_GetSettings(t *testing.T) {
	type fields struct {
		number   int64
		settings *setup.Settings
	}
	tests := []struct {
		name    string
		fields  fields
		want    *setup.Settings
		wantErr bool
	}{
		// test cases
		{
			name:    "simple",
			fields:  fields{settings: &setup.Settings{Step: 1, Limit: 10}, number: 1},
			want:    &setup.Settings{Step: 1, Limit: 10},
			wantErr: false,
		},
		{
			name:    "not simple",
			fields:  fields{settings: &setup.Settings{Step: 2, Limit: 20}, number: 1},
			want:    &setup.Settings{Step: 2, Limit: 20},
			wantErr: false,
		},
		{
			name:    "error: constructor required",
			fields:  fields{settings: nil, number: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		c := &Counter{
			number:   tt.fields.number,
			settings: tt.fields.settings,
		}
		got, err := c.GetSettings()
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Counter.GetSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Counter.GetSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_checkSettings(t *testing.T) {
	type args struct {
		s *setup.Settings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "correct settings",
			args:    args{s: &setup.Settings{Step: 2, Limit: 20}},
			wantErr: false,
		},
		{
			name:    "incorrect settings",
			args:    args{s: &setup.Settings{Step: 20, Limit: 2}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if err := checkSettings(tt.args.s); (err != nil) != tt.wantErr {
			t.Errorf("%q. checkSettings() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
