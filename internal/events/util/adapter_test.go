package util

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseTime(t *testing.T) {
	type args struct {
		timestamp int64
	}

	now := time.Now()
	secs := now.Unix()
	msecs := now.UnixMilli()
	usecs := now.UnixMicro()
	nsecs := now.UnixNano()

	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Success Seconds",
			args:    args{secs},
			want:    now.Round(time.Second),
			wantErr: false,
		},
		{
			name:    "Success MilliSeconds",
			args:    args{msecs},
			want:    now.Round(time.Millisecond),
			wantErr: false,
		},
		{
			name:    "Success MicroSeconds",
			args:    args{usecs},
			want:    now.Round(time.Microsecond),
			wantErr: false,
		},
		{
			name:    "Success NanoSeconds",
			args:    args{nsecs},
			want:    now.Round(time.Nanosecond),
			wantErr: false,
		},
		{
			name:    "Negative value error",
			args:    args{-12323131312},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Not enough precision error",
			args:    args{123123},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Timestamp: %d", tt.args.timestamp)
			got, err := parseTime(tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
