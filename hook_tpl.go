// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/9

package fsconf

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
)

var _ Hook = (*HookTPL)(nil)

// HookTPL 一个通用的可用于替换内容的 Hook 模版
type HookTPL struct {
	// HookName 名称，必填
	HookName string

	// KeyPrefix 查找的表达式前缀，必填
	// 最终表达是为  {$RegexpKey.([A-Za-z0-9_]+)}
	KeyPrefix string

	// Values 用于查找的值，必填
	Values map[string]string

	regexp *regexp.Regexp

	initOnce sync.Once
	initErr  error
}

func (h *HookTPL) Name() string {
	return h.HookName
}

func (h *HookTPL) init() {
	if h.HookName == "" {
		h.initErr = errors.New("empty HookName")
		return
	}
	if h.KeyPrefix == "" {
		h.initErr = errors.New("empty KeyPrefix")
		return
	}
	rule := fmt.Sprintf(`\{%s\.([A-Za-z0-9_]+)\}`, h.KeyPrefix)
	h.regexp, h.initErr = regexp.Compile(rule)
}

func (h *HookTPL) Execute(ctx context.Context, p *HookParam) (output []byte, err error) {
	h.initOnce.Do(h.init)
	if h.initErr != nil {
		return nil, h.initErr
	}

	contentNew := h.regexp.ReplaceAllFunc(p.Content, func(subStr []byte) []byte {
		// 将 {prefix.xxx} 中的 xxx 部分取出
		key := subStr[len("{"+h.KeyPrefix+".") : len(subStr)-1] // eg: xxx
		var val string
		val, err = h.getValue(string(key))
		if err != nil {
			return nil
		}
		return []byte(val)
	})
	if err != nil {
		return nil, err
	}
	return contentNew, err
}

func (h *HookTPL) getValue(key string) (string, error) {
	if len(h.Values) == 0 {
		return "", fmt.Errorf("empty Values for hook %q", h.Name())
	}
	val, ok := h.Values[key]
	if !ok {
		return "", fmt.Errorf("key=%q not found for hook %q", key, h.Name())
	}
	return val, nil
}
