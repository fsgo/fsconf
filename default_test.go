/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"reflect"
	"testing"

	"github.com/fsgo/fsenv"
)

func init() {
	env := fsenv.NewAppEnv(&fsenv.Value{RootDir: "./testdata"})
	Default.SetEnvOnce(env)
}

func TestExists(t *testing.T) {
	type args struct {
		confName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				confName: "not_exists.toml",
			},
			want: false,
		},
		{
			name: "case 2",
			args: args{
				confName: "abc.json",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.confName); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		confName string
		obj      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				confName: "not_exists.json",
				obj:      map[string]string{},
			},
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				confName: "abc.json",
				obj:      map[string]string{},
			},
			want: map[string]string{
				"A": "bb",
			},
			wantErr: false,
		},
		{
			name: "case 3",
			args: args{
				confName: "db1.toml",
				obj:      map[string]string{},
			},
			want: map[string]string{
				"name":    "abc",
				"charset": "utf-8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Parse(tt.args.confName, &tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.args.obj
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestParseByAbsPath(t *testing.T) {
	type args struct {
		confAbsPath string
		obj         map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				confAbsPath: "testdata/conf/abc.json",
				obj:         map[string]string{},
			},
			wantErr: false,
			want: map[string]string{
				"A": "bb",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseByAbsPath(tt.args.confAbsPath, &tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("ParseByAbsPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.args.obj
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseByAbsPath() got=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestRegisterHelper(t *testing.T) {
	type args struct {
		name string
		fn   HelperFn
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				name: "",
				fn:   helperOsEnvVars,
			},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				name: "test_helper",
				fn:   helperOsEnvVars,
			},
			wantErr: false,
		},
		{
			name: "case 3- name is same as case 2",
			args: args{
				name: "test_helper",
				fn:   helperOsEnvVars,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterHelper(tt.args.name, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("RegisterHelper() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterParser(t *testing.T) {
	type args struct {
		fileExt string
		fn      ParserFn
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				fileExt: ".json",
				fn:      JSONParser,
			},
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				fileExt: ".myjson",
				fn:      JSONParser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterParser(tt.args.fileExt, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("RegisterParser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
