/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/4
 */

package parser

import (
	"fmt"
	"reflect"
	"testing"
)

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
				input: []byte(`line1
#line2
  #line3
#line4

line6 #666`),
			},
			wantOut: []byte(`line1




line6 #666`),
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

func TestGetDefault(t *testing.T) {
	type args struct {
		ext string
	}
	tests := []struct {
		name string
		args args
		want Fn
	}{
		{
			name: "case 1",
			args: args{
				ext: FileJSON,
			},
			want: JSON,
		},
		{
			name: "case 2",
			args: args{
				ext: ".other_not_found",
			},
			want: nil,
		},
		{
			name: "case 3",
			args: args{
				ext: FileTOML,
			},
			want: TOML,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDefault(tt.args.ext)
			if got == nil && tt.want == nil {
				return
			}
			gotName := fmt.Sprint(got)
			wantName := fmt.Sprint(tt.want)
			if gotName != wantName {
				t.Errorf("GetDefault() = %v, want %v", gotName, wantName)
			}
		})
	}
}
