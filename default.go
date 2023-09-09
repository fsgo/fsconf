// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"context"
	"sync/atomic"
)

var defaultCfg atomic.Pointer[Configure]

func init() {
	defaultCfg.Store(NewDefault())
}

// Default 默认的实例
func Default() *Configure {
	return defaultCfg.Load()
}

func SetDefault(cfg *Configure) (old *Configure) {
	return defaultCfg.Swap(cfg)
}

// Parse 解析配置，配置文件默认认为在 conf/目录下,
// 如 有 conf/abc.toml ，则 confName="abc.toml"
func Parse(confName string, obj any) (err error) {
	return Default().Parse(confName, obj)
}

// ParseByAbsPath 解析绝对路径的配置
func ParseByAbsPath(confAbsPath string, obj any) (err error) {
	return Default().ParseByAbsPath(confAbsPath, obj)
}

// ParseBytes （全局）解析 bytes
// fileExt 是文件后缀，如.json、.toml
func ParseBytes(fileExt string, content []byte, obj any) error {
	return Default().ParseBytes(fileExt, content, obj)
}

// Exists  （全局）判断是否存在
func Exists(confName string) bool {
	return Default().Exists(confName)
}

// RegisterParser （全局）注册一个解析器
// fileExt 是文件后缀，如 .json
func RegisterParser(fileExt string, fn ParserFn) error {
	err := Default().RegisterParser(fileExt, fn)
	if err != nil {
		return err
	}
	defaultParsers = append(defaultParsers, parserNameFn{Name: fileExt, Fn: fn})
	return nil
}

// RegisterHook （全局）注册一个辅助类
func RegisterHook(h Hook) error {
	if err := defaultHooks.Add(h); err != nil {
		return err
	}
	return Default().RegisterHook(h)
}

// MustRegisterHook （全局）注册一个辅助类，若失败会 panic
func MustRegisterHook(h Hook) {
	if err := RegisterHook(h); err != nil {
		panic(err)
	}
}

// WithContext （全局）返回新的对象,并设置新的 ctx
func WithContext(ctx context.Context) *Configure {
	return Default().WithContext(ctx)
}

// WithHook （全局）返回新的对象,并注册 Hook
func WithHook(hs ...Hook) *Configure {
	return Default().WithHook(hs...)
}
