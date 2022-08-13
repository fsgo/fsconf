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
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/fsgo/fsconf/internal/parser"
)

var _ Hook = (*hookInclude)(nil)

type hookInclude struct{}

func (h *hookInclude) Name() string {
	return "template"
}

var hookTplEnableReg = regexp.MustCompile(`hook\.template\s+Enable=true`)

func (h *hookInclude) Execute(ctx context.Context, p *HookParam) (output []byte, err error) {
	cmts := parser.HeadComments(p.Content)
	if len(cmts) == 0 {
		return p.Content, nil
	}
	var enable bool
	for _, cmt := range cmts {
		if hookTplEnableReg.MatchString(cmt) {
			enable = true
			break
		}
	}
	if !enable {
		return p.Content, nil
	}
	return h.exec(ctx, p)
}

func (h *hookInclude) exec(ctx context.Context, p *HookParam) (output []byte, err error) {
	tmpl := template.New("config")
	tmpl.Delims("{template", "template}")
	tmpl.Funcs(map[string]any{
		"include": func(name string) (string, error) {
			return h.fnInclude(ctx, name, p)
		},
	})
	tmpl, err = tmpl.Parse(string(p.Content))
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err = tmpl.Execute(buf, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *hookInclude) pathHasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

func (h *hookInclude) fnInclude(ctx context.Context, name string, p *HookParam) (string, error) {
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
		o1, err2 := h.exec(ctx, p1)
		if err2 != nil {
			return "", err2
		}
		buf.Write(o1)
	}
	return buf.String(), nil
}
