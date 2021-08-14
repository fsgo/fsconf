// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package parser

// Fn 对应文件后缀的配置解析方法
type Fn func(bf []byte, obj interface{}) error

const (
	// 已支持的文件后缀

	// FileTOML toml
	FileTOML = ".toml"

	// FileJSON  json
	FileJSON = ".json"
)
