/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
)

// ParserFn 对应文件后缀的配置解析方法
type ParserFn func(bf []byte, obj interface{}) error

const (
	// 已支持的文件后缀

	// FileExtTOML toml
	FileExtTOML = ".toml"

	// FileExtJSON  json
	FileExtJSON = ".json"
)

// JSONParser .json 文件的解析方法
// 若内容以 # 开头，则该为注释
func JSONParser(txt []byte, obj interface{}) error {
	bf := StripComment(txt)
	dec := json.NewDecoder(bytes.NewReader(bf))
	dec.UseNumber()
	return dec.Decode(obj)
}

// StripComment 去除注释
func StripComment(input []byte) (out []byte) {
	var buf bytes.Buffer
	lines := bytes.Split(input, []byte("\n"))
	for _, line := range lines {
		if !bytes.HasPrefix(line, []byte("#")) {
			buf.Write(line)
			buf.WriteString("\n")
		}
	}
	return bytes.TrimSpace(buf.Bytes())
}

var _ ParserFn = JSONParser

var defaultParsers = map[string]ParserFn{
	FileExtJSON: JSONParser,
	FileExtTOML: toml.Unmarshal,
}
