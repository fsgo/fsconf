// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"context"
)

// Default 默认的实例
var Default = NewDefault()

// Parse 解析配置，配置文件默认认为在 conf/目录下,
// 如 有 conf/abc.toml ，则 confName="abc.toml"
func Parse(confName string, obj any) (err error) {
	return Default.Parse(confName, obj)
}

// ParseByAbsPath 解析绝对路径的配置
func ParseByAbsPath(confAbsPath string, obj any) (err error) {
	return Default.ParseByAbsPath(confAbsPath, obj)
}

// ParseBytes （全局）解析 bytes
// fileExt 是文件后缀，如.json、.toml
func ParseBytes(fileExt string, content []byte, obj any) error {
	return Default.ParseBytes(fileExt, content, obj)
}

// Exists  （全局）判断是否存在
func Exists(confName string) bool {
	return Default.Exists(confName)
}

// RegisterParser （全局）注册一个解析器
// fileExt 是文件后缀，如 .json
func RegisterParser(fileExt string, fn ParserFn) error {
	defaultParsers[fileExt] = fn
	return Default.RegisterParser(fileExt, fn)
}

// RegisterHelper （全局）注册一个辅助方法
func RegisterHelper(h Hook) error {
	_ = defaultHooks.Add(h)
	return Default.RegisterHook(h)
}

// WithContext （全局）设置一个 context，并返回新的对象
func WithContext(ctx context.Context) Configure {
	return Default.WithContext(ctx)
}
