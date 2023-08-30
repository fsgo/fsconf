// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"context"
	"errors"
	"fmt"

	"github.com/fsgo/fsconf/internal/hook"
)

// Hook 辅助方法，在执行解析前，会先会配置的内容进行解析处理
type Hook interface {
	Name() string
	Execute(ctx context.Context, p *HookParam) (output []byte, err error)
}

// HookParam param for helper
type HookParam struct {
	FileExt   string    // 文件类型后缀，如 .toml,.json
	Configure Configure // 当前 Configure 对象
	ConfPath  string    // 文件路径。当直接解析内容时，为空字符串
	Content   []byte    // 文件内容
}

var defaultHooks hooks = []Hook{
	&hookTemplate{},
	newHook("osenv", hook.OsEnvVars),
	&hookFsEnv{},
}

type hooks []Hook

func (hs *hooks) Add(h Hook) error {
	if len(h.Name()) == 0 {
		return errors.New("hook.Name is empty, not allow")
	}

	for _, h1 := range *hs {
		if h.Name() == h1.Name() {
			return fmt.Errorf("hook=%q already exists", h.Name())
		}
	}
	*hs = append(*hs, h)
	return nil
}

func (hs hooks) Execute(ctx context.Context, p *HookParam) (output []byte, err error) {
	if len(hs) == 0 {
		return p.Content, nil
	}
	input := p.Content
	content := make([]byte, len(input))
	if n := copy(content, input); n != len(input) {
		return nil, fmt.Errorf("copy config content failed, want=%d copied=%d", len(input), n)
	}

	for _, hk := range hs {
		p.Content = content
		content, err = hk.Execute(ctx, p)
		if err != nil {
			return nil, fmt.Errorf("hook=%q has error:%w", hk.Name(), err)
		}
	}
	return content, err
}

type hookTpl struct {
	fn   hook.Fn
	name string
}

func (h *hookTpl) Name() string {
	return h.name
}

func (h *hookTpl) Execute(_ context.Context, p *HookParam) (output []byte, err error) {
	return h.fn(p.ConfPath, p.Content)
}

func newHook(name string, fn hook.Fn) Hook {
	return &hookTpl{
		name: name,
		fn:   fn,
	}
}
