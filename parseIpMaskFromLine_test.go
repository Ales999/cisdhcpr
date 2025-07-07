package main

import (
	"net/netip"
	"reflect"
	"testing"
)

func Test_parseIpMaskFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    netip.Addr
		want1   netip.Prefix
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test case 1",
			args: args{
				line: " ip address 172.24.62.201 255.255.255.248",
			},
			want:    netip.MustParseAddr("172.24.62.201"),
			want1:   netip.MustParsePrefix("172.24.62.201/29"),
			wantErr: false,
		},
		{
			name: "Test case 2",
			args: args{
				line: " ip address 10.0.0.1 255.255.255.0",
			},
			want:    netip.MustParseAddr("10.0.0.1"),
			want1:   netip.MustParsePrefix("10.0.0.1/24"),
			wantErr: false,
		},
		{
			name: "Test case 3",
			args: args{
				line: " ip address 192.168.0.1 255.255.0.0",
			},
			want:    netip.MustParseAddr("192.168.0.1"),
			want1:   netip.MustParsePrefix("192.168.0.1/16"),
			wantErr: false,
		},
		{
			name: "Test case 4",
			args: args{
				line: " ip address 172.16.0.1 255.255.240.0",
			},
			want:    netip.MustParseAddr("172.16.0.1"),
			want1:   netip.MustParsePrefix("172.16.0.1/20"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseIpMaskFromLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIpMaskFromLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIpMaskFromLine() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseIpMaskFromLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
