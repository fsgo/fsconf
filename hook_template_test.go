// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/15

package fsconf

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func Test_hookInclude_Execute(t *testing.T) {
	getHookParam := func(fp string) *HookParam {
		bf, err := os.ReadFile(fp)
		if err != nil {
			panic(err)
		}
		return &HookParam{
			ConfPath: fp,
			Content:  bf,
		}
	}

	type args struct {
		ctx  context.Context
		getP func() *HookParam
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []byte
		wantErr    bool
	}{
		{
			name: "include.toml",
			args: args{
				ctx: context.Background(),
				getP: func() *HookParam {
					return getHookParam("testdata/conf/include.toml")
				},
			},
			wantOutput: []byte("# hook.template  Enable=true\nA=\"a\"\n\nB=\"b\"\nB1=\"b1\"\nC=\"c\"\n\nZ=\"z\"\n"),
		},
		{
			name: "include not found",
			args: args{
				ctx: context.Background(),
				getP: func() *HookParam {
					return getHookParam("testdata/conf/include_e1.toml")
				},
			},
			wantErr: true,
		},
		{
			name: "include ConfPath Empty",
			args: args{
				ctx: context.Background(),
				getP: func() *HookParam {
					p := getHookParam("testdata/conf/include.toml")
					p.ConfPath = ""
					return p
				},
			},
			wantErr: true,
		},
		{
			name: "include not enable",
			args: args{
				ctx: context.Background(),
				getP: func() *HookParam {
					p := getHookParam("testdata/conf/include_not_enable.toml")
					p.ConfPath = ""
					return p
				},
			},
			wantOutput: []byte("A=\"a\"\n\n{template include \"not_found.toml\" template}\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hookInclude{}
			gotOutput, err := h.Execute(tt.args.ctx, tt.args.getP())
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Execute() gotOutput = %q, want %q", gotOutput, tt.wantOutput)
			}
		})
	}
}
