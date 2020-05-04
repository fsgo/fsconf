/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/4
 */

package helper

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		input   []byte
		helpers []*Helper
	}
	helpers := []*Helper{
		New("no_d", func(confContent []byte) ([]byte, error) {
			if bytes.Contains(confContent, []byte("error")) {
				return nil, fmt.Errorf("must error")
			}
			return bytes.ReplaceAll(confContent, []byte("d"), []byte("")), nil
		}),
	}

	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{
			name: "case 1",
			args: args{
				input:   []byte("abcd"),
				helpers: nil,
			},
			wantOutput: []byte("abcd"),
			wantErr:    false,
		},
		{
			name: "case 2",
			args: args{
				input:   []byte("abcd"),
				helpers: helpers,
			},
			wantOutput: []byte("abc"),
			wantErr:    false,
		},
		{
			name: "case 3-error",
			args: args{
				input:   []byte("abcd has error"),
				helpers: helpers,
			},
			wantOutput: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := Execute(tt.args.input, tt.args.helpers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Execute() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
