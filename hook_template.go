// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/14

package fsconf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/fsgo/fsenv"

	"github.com/fsgo/fsconf/internal/parser"
	"github.com/fsgo/fsconf/internal/xcache"
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
	params := make(map[string]string, 3)
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
		"fetch": func(name string, args ...string) (string, error) {
			return h.fnFetch(ctx, hp, tp, name, args)
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

	data := map[string]string{
		"IDC":         fsenv.IDC(),
		"RootDir":     fsenv.RootDir(),
		"ConfRootDir": fsenv.ConfDir(),
		"LogRootDir":  fsenv.LogDir(),
		"DataRootDir": fsenv.DataDir(),
		"RunMode":     fsenv.RunMode().String(),
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
	if len(p.ConfPath) == 0 {
		return "", errors.New("p.ConfPath is empty cannot use include")
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

func (h *hookTemplate) getXCache(p *HookParam) *xcache.FileCache {
	dir := filepath.Join(fsenv.TempDir(), "fsconf_cache")
	return &xcache.FileCache{
		Dir: dir,
	}
}

func (h *hookTemplate) fnFetch(ctx context.Context, p *HookParam, tp map[string]string, api string, ps []string) (string, error) {
	if len(api) == 0 {
		return "", errors.New("url is required")
	}
	if len(ps) > 1 {
		return "", errors.New("only support 0 or 1 param")
	}

	timeout := 3 * time.Second
	var cacheTTL time.Duration
	if len(ps) == 1 {
		param, err := xcache.ParserParam(ps[0])
		if err != nil {
			return "", err
		}
		if param.Timeout > 0 {
			timeout = param.Timeout
		}
		cacheTTL = param.TTL
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	bf, err := httpFetch(ctx, api)

	if cacheTTL > 0 {
		fc := h.getXCache(p)
		if err == nil {
			fc.Set(api, bf)
		} else {
			if cv, ok := fc.Get(api, cacheTTL); ok {
				return string(cv), nil
			}
		}
	}

	return string(bf), err
}

func httpFetch(ctx context.Context, api string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
