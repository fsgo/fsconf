// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package parser

import (
	"bytes"
	"encoding/json"
)

// JSON .json 文件的解析方法
// 若内容以 # 开头，则该为注释
func JSON(txt []byte, obj any) error {
	bf := StripComment(txt)
	dec := json.NewDecoder(bytes.NewReader(bf))
	dec.UseNumber()
	return dec.Decode(obj)
}
