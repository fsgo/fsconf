// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"encoding/xml"

	"github.com/fsgo/fsconf/internal/parser"
)

// DecoderFunc 针对特定文件后缀的配置解析方法
// 当前已经内置了 .toml  和 .json的解析方法
type DecoderFunc func(bf []byte, obj any) error

type parserNameFn struct {
	Fn   DecoderFunc
	Name string
}

// defaultParsers 所有默认的 parser，
// 当传入配置文件名不包含后置的时候，会使用此顺序依次查找
var defaultParsers = []parserNameFn{
	{Name: ".json", Fn: parser.JSON},
	{Name: ".xml", Fn: xml.Unmarshal},
}
