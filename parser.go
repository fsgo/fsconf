// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/7/11

package fsconf

import (
	"encoding/xml"

	"gopkg.in/yaml.v2"

	"github.com/fsgo/fsconf/internal/parser"
)

// ParserFn 针对特定文件后缀的配置解析方法
// 当前已经内置了 .toml  和 .json的解析方法
type ParserFn = parser.Fn

// defaultParsers 所有默认的parser
var defaultParsers = map[string]parser.Fn{
	parser.FileJSON: parser.JSON,
	parser.FileTOML: parser.TOML,
	".xml":          xml.Unmarshal,
	".yml":          yaml.Unmarshal,
}
