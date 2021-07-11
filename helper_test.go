// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestHelpersExecute(t *testing.T) {

	var hs helpers
	_ = hs.Add(newHelper("no_d", func(confContent []byte) ([]byte, error) {
		if bytes.Contains(confContent, []byte("error")) {
			return nil, fmt.Errorf("must error")
		}
		return bytes.ReplaceAll(confContent, []byte("d"), []byte("")), nil
	}))

	if len(hs) == 0 {
		t.Fatal("helper is empty")
	}

	type args struct {
		input   []byte
		helpers helpers
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
				helpers: hs,
			},
			wantOutput: []byte("abc"),
			wantErr:    false,
		},
		{
			name: "case 3-error",
			args: args{
				input:   []byte("abcd has error"),
				helpers: hs,
			},
			wantOutput: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := tt.args.helpers.Execute(context.Background(), nil, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("WithFunc() gotOutput = %q, want %q", gotOutput, tt.wantOutput)
			}
		})
	}
}
