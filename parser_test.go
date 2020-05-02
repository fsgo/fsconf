/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"reflect"
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
			if err := JSONParser(tt.args.txt, &tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("JSONParser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripComment(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name    string
		args    args
		wantOut []byte
	}{
		{
			name: "case 1",
			args: args{
				input: []byte(`abc
#你好
  #abc
#abd
def`),
			},
			wantOut: []byte(`abc
  #abc
def`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := StripComment(tt.args.input); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("StripComment() = %q, want %q", gotOut, tt.wantOut)
			}
		})
	}
}
