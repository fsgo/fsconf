// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/8

package fsconf

import (
	"context"
	"reflect"
	"testing"
)

func Test_fsEnvHelper_getValue(t *testing.T) {
	type args struct {
		key string
		cf  *Configure
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "RootDir",
			args: args{
				key: "RootDir",
				cf:  Default(),
			},
			want: "./testdata",
		},
		{
			name: "IDC",
			args: args{
				key: "IDC",
				cf:  Default(),
			},
			want: "test",
		},
		{
			name: "DataRootDir",
			args: args{
				key: "DataRootDir",
				cf:  Default(),
			},
			want: "testdata/data",
		},
		{
			name: "ConfRootDir",
			args: args{
				key: "ConfRootDir",
				cf:  Default(),
			},
			want: "testdata/conf",
		},
		{
			name: "LogRootDir",
			args: args{
				key: "LogRootDir",
				cf:  Default(),
			},
			want: "testdata/log",
		},
		{
			name: "RunMode",
			args: args{
				key: "RunMode",
				cf:  Default(),
			},
			want: "product",
		},
		{
			name: "other key not support",
			args: args{
				key: "other-key",
				cf:  Default(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &hookFsEnv{}
			got, err := f.getValue(tt.args.key, tt.args.cf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fsEnvHelper_Execute(t *testing.T) {
	type args struct {
		ctx   context.Context
		cf    *Configure
		input []byte
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{
			name: "idc and log dir",
			args: args{
				cf:    Default(),
				ctx:   context.Background(),
				input: []byte(`{"idc":"{fsenv.IDC}","logDir":"{fsenv.LogRootDir}"}`),
			},
			wantOutput: []byte(`{"idc":"test","logDir":"testdata/log"}`),
		},
		{
			name: "not support key",
			args: args{
				cf:    Default(),
				ctx:   context.Background(),
				input: []byte(`{"idc":"{fsenv.other}"}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &hookFsEnv{}
			p := &HookParam{
				Configure: tt.args.cf,
				Content:   tt.args.input,
			}
			gotOutput, err := f.Execute(tt.args.ctx, p)
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
