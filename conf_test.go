// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"reflect"
	"testing"

	"github.com/fsgo/fsenv"

	"github.com/fsgo/fsconf/internal/helper"
	"github.com/fsgo/fsconf/internal/parser"
)

func Test_confImpl(t *testing.T) {
	conf := New()
	env := fsenv.NewAppEnv(fsenv.Value{
		RootDir: "./testdata",
	})
	conf.SetEnvOnce(env)
	var a interface{}
	if err := conf.Parse("abc.json", &a); err == nil {
		t.Errorf("expect has error")
	}

	conf.RegisterParser(parser.FileJSON, parser.JSON)

	if err := conf.Parse("abc.xyz", &a); err == nil {
		t.Errorf("expect has error 2")
	}
}

func TestNewDefault1(t *testing.T) {
	hd := append([]*helper.Helper{}, helper.Defaults...)
	defer func() {
		helper.Defaults = hd
		if re := recover(); re == nil {
			t.Errorf("want panic")
		}
	}()
	h := helper.New("test", helper.OsEnvVars)
	// helper 有重复的时候
	helper.Defaults = append(helper.Defaults, h, h)
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
