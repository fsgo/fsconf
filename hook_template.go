// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/14

package fsconf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/fsgo/fsenv"

	"github.com/fsgo/fsconf/internal/parser"
)

var _ Hook = (*hookTemplate)(nil)

type hookTemplate struct{}

func (h *hookTemplate) Name() string {
	return "template"
}

var hookTplPrefix = "hook.template "

func (h *hookTemplate) Execute(ctx context.Context, hp *HookParam) (output []byte, err error) {
	cmts := parser.HeadComments(hp.Content)
	if len(cmts) == 0 {
		return hp.Content, nil
	}
	params := make(map[string]string)
	for _, cmt := range cmts {
		if strings.HasPrefix(cmt, hookTplPrefix) {
			arr := strings.Fields(cmt[len(hookTplPrefix):])
			for i := 0; i < len(arr); i++ {
				tmp := strings.Split(arr[i], "=")
				if len(tmp) == 2 && len(tmp[0]) > 0 && len(tmp[1]) > 0 {
					params[tmp[0]] = tmp[1]
				}
			}
		}
	}
	if params["Enable"] != "true" {
		return hp.Content, nil
	}
	return h.exec(ctx, hp, params)
}

func (h *hookTemplate) exec(ctx context.Context, hp *HookParam, tp map[string]string) (output []byte, err error) {
	tmpl := template.New("config")
	left := "{{"
	right := "}}"
	if v := tp["Left"]; len(v) > 0 {
		left = v
	}
	if v := tp["Right"]; len(v) > 0 {
		right = v
	}
	tmpl.Delims(left, right)
	tmpl.Funcs(map[string]any{
		"include": func(name string) (string, error) {
			return h.fnInclude(ctx, name, hp, tp)
		},
		"osenv": func(name string) string {
			return os.Getenv(name)
		},
		"contains": func(s string, sub string) bool {
			return strings.Contains(s, sub)
		},
		"prefix": func(s string, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"suffix": func(s string, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
	})
	tmpl, err = tmpl.Parse(string(hp.Content))
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	data := make(map[string]string, 6)
	if cae, ok := hp.Configure.(fsenv.HasAppEnv); ok {
		ce := cae.AppEnv()
		data["IDC"] = ce.IDC()
		data["RootDir"] = ce.RootDir()
		data["ConfRootDir"] = ce.ConfRootDir()
		data["LogRootDir"] = ce.LogRootDir()
		data["DataRootDir"] = ce.DataRootDir()
		data["RunMode"] = string(ce.RunMode())
	}

	if err = tmpl.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *hookTemplate) pathHasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

func (h *hookTemplate) fnInclude(ctx context.Context, name string, p *HookParam, tp map[string]string) (string, error) {
	if p.ConfPath == "" {
		return "", fmt.Errorf("p.ConfPath is empty cannot use include")
	}
	var fp string
	if filepath.IsAbs(name) {
		fp = name
	} else {
		fp = filepath.Join(filepath.Dir(p.ConfPath), name)
	}

	files, err := filepath.Glob(fp)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		if !h.pathHasMeta(name) {
			return "", fmt.Errorf("include %q not found", name)
		}
		return "", nil
	}
	var buf bytes.Buffer
	for _, f := range files {
		body, err1 := os.ReadFile(f)
		if err1 != nil {
			return "", err1
		}

		p1 := &HookParam{
			ConfPath:  f,
			Content:   body,
			Configure: p.Configure,
		}
		o1, err2 := h.exec(ctx, p1, tp)
		if err2 != nil {
			return "", err2
		}
		buf.Write(o1)
	}
	return buf.String(), nil
}