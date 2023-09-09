// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/15

package fsconf

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_hookInclude_Execute(t *testing.T) {
	getHookParam := func(fp string) *HookParam {
		bf, err := os.ReadFile(fp)
		if err != nil {
			panic(err)
		}
		return &HookParam{
			ConfPath:  fp,
			Content:   bf,
			Configure: Default,
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
			wantOutput: []byte(
				`# hook.template  Enable=true
A="a"
Port = {osenv.Port1}

B="b"
B1="b1"
C="c"

Z="z"


`),
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
			wantOutput: []byte("A=\"a\"\n\n{{ include \"not_found.toml\" }}\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hookTemplate{}
			gotOutput, err := h.Execute(tt.args.ctx, tt.args.getP())
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, string(tt.wantOutput), string(gotOutput))
		})
	}
}

func Test_fnFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k := r.URL.Query().Get("k")
		_, _ = w.Write([]byte("hello-" + k))
	}))
	defer ts.Close()
	api := ts.URL
	t.Run("server ok", func(t *testing.T) {
		txt := `# hook.template  Enable=true
K1="{{ fetch "` + api + `?k=k1" }}"
K2="{{ fetch "` + api + `?k=k2" }}"
K3="{{ fetch "` + api + `?k=k3" "timeout=5s&cache=1h" }}"
`
		mp := map[string]string{}
		err1 := ParseBytes(".toml", []byte(txt), &mp)
		require.NoError(t, err1)
		want1 := map[string]string{
			"K1": "hello-k1",
			"K2": "hello-k2",
			"K3": "hello-k3",
		}
		require.Equal(t, want1, mp)
	})

	t.Run("server unreachable", func(t *testing.T) {
		ts.Close()
		txt := `# hook.template  Enable=true
K3="{{ fetch "` + api + `?k=k3" "timeout=5s&cache=1h" }}"
`
		mp := map[string]string{}
		err1 := ParseBytes(".toml", []byte(txt), &mp)
		require.NoError(t, err1)
		want1 := map[string]string{
			"K3": "hello-k3",
		}
		require.Equal(t, want1, mp)
	})
}
