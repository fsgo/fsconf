/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/4
 */

package parser

import (
	"bytes"
)

// Fn 对应文件后缀的配置解析方法
type Fn func(bf []byte, obj interface{}) error

const (
	// 已支持的文件后缀

	// FileTOML toml
	FileTOML = ".toml"

	// FileJSON  json
	FileJSON = ".json"
)

// StripComment 去除单行的'#'注释
// 只支持单行，不支持行尾
func StripComment(input []byte) (out []byte) {
	var buf bytes.Buffer
	lines := bytes.Split(input, []byte("\n"))
	for _, line := range lines {
		lineN := bytes.TrimSpace(line)
		if !bytes.HasPrefix(lineN, []byte("#")) {
			buf.Write(line)
		}
		buf.WriteString("\n")
	}
	return bytes.TrimSpace(buf.Bytes())
}

// Defaults 所有默认的parser
var Defaults = map[string]Fn{
	FileJSON: JSON,
	FileTOML: TOML,
}

// GetDefault 获取指定默认的parser
func GetDefault(ext string) Fn {
	fn, _ := Defaults[ext]
	return fn
}
