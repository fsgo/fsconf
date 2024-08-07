// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"reflect"
	"testing"

	"github.com/fsgo/fst"

	"github.com/fsgo/fsconf/internal/hook"
	"github.com/fsgo/fsconf/internal/parser"
)

func Test_confImpl(t *testing.T) {
	conf := New()
	testReset()
	var a any
	fst.Error(t, conf.Parse("abc.json", &a))
	fst.NoError(t, conf.RegisterParser(".json", parser.JSON))
	fst.Error(t, conf.Parse("abc.xyz", &a))
	fst.NoError(t, conf.Parse("testdata/db10.json", &a))
}

func TestNewDefault1(t *testing.T) {
	hd := append([]Hook{}, defaultHooks...)
	defer func() {
		defaultHooks = hd
		if re := recover(); re == nil {
			t.Errorf("want panic")
		}
	}()
	h := newHook("test", hook.OsEnvVars)
	// helper 有重复的时候
	defaultHooks = append(defaultHooks, h, h)
	NewDefault()
}

func Test_confImpl_ParseBytes(t *testing.T) {
	type args struct {
		fileExt string
		content []byte
		obj     map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    map[string]string
	}{
		{
			name: "case 1",
			args: args{
				fileExt: "",
				content: nil,
				obj:     map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				fileExt: ".json",
				content: []byte(`{"Name":"Hello"}`),
				obj:     map[string]string{},
			},
			wantErr: false,
			want:    map[string]string{"Name": "Hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewDefault()
			if err := c.ParseBytes(tt.args.fileExt, tt.args.content, &tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("ParseBytes() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(tt.args.obj, tt.want) {
					t.Errorf("ParseBytes(), obj=%v, got=%v", tt.args.obj, tt.want)
				}
			}
		})
	}
}
