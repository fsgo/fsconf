// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"context"
	"fmt"

	"github.com/fsgo/fsconf/internal/helper"
)

// Helper 辅助方法，在执行解析前，会先会配置的内容进行解析处理
type Helper interface {
	Name() string
	Execute(ctx context.Context, cf Configure, input []byte) (output []byte, err error)
}

var defaultHelpers helpers = []Helper{
	newHelper("osenv", helper.OsEnvVars),
}

type helpers []Helper

func (hs *helpers) Add(h Helper) error {
	if h.Name() == "" {
		return fmt.Errorf("helper.Name is empty, not allow")
	}

	for _, h1 := range *hs {
		if h.Name() == h1.Name() {
			return fmt.Errorf("helper=%q already exists", h.Name())
		}
	}
	*hs = append(*hs, h)
	return nil
}

func (hs helpers) Execute(ctx context.Context, cf Configure, input []byte) (output []byte, err error) {
	if len(hs) == 0 {
		return input, nil
	}
	content := make([]byte, len(input))
	if n := copy(content, input); n != len(input) {
		return nil, fmt.Errorf("copy failed")
	}

	for _, helper := range hs {
		content, err = helper.Execute(ctx, cf, content)
		if err != nil {
			return nil, fmt.Errorf("helper=%q has error:%w", helper.Name(), err)
		}
	}
	return content, err
}

// Helper 辅助功能
// 在正式解析配置前执行
type helperTpl struct {
	name string
	fn   helper.Fn
}

func (h *helperTpl) Name() string {
	return h.name
}

func (h *helperTpl) Execute(ctx context.Context, cf Configure, input []byte) (output []byte, err error) {
	return h.fn(input)
}

// New 创建实例
func newHelper(name string, fn helper.Fn) Helper {
	return &helperTpl{
		name: name,
		fn:   fn,
	}
}
