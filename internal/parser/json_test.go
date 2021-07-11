// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package parser

import (
	"testing"
)

func Test_JSONParser(t *testing.T) {
	type args struct {
		txt []byte
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				txt: []byte(""),
				obj: nil,
			},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				txt: []byte("abc"),
				obj: nil,
			},
			wantErr: true,
		},
		{
			name: "case 3",
			args: args{
				txt: []byte(`{"a":"b"}`),
				obj: map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := JSON(tt.args.txt, &tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("JSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
