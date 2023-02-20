// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"encoding/xml"

	"gopkg.in/yaml.v3"

	"github.com/fsgo/fsconf/internal/parser"
)

// ParserFn 针对特定文件后缀的配置解析方法
// 当前已经内置了 .toml  和 .json的解析方法
type ParserFn = parser.Fn

type parserNameFn struct {
	Fn   ParserFn
	Name string
}

// defaultParsers 所有默认的 parser，
// 当传入配置文件名不包含后置的时候，会使用此顺序依次查找
var defaultParsers = []parserNameFn{
	{Name: parser.FileTOML, Fn: parser.TOML},
	{Name: parser.FileJSON, Fn: parser.JSON},
	{Name: ".yml", Fn: yaml.Unmarshal},
	{Name: ".xml", Fn: xml.Unmarshal},
}
