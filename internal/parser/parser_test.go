// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package parser

import (
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
