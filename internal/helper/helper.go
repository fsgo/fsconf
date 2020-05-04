/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/4
 */

package helper

import (
	"fmt"
)

// Fn helper 的函数
type Fn func(confContent []byte) ([]byte, error)

// Helper 辅助功能
// 在正式解析配置前执行
type Helper struct {
	Name string
	Fn   Fn
}

// New 创建实例
func New(name string, fn Fn) *Helper {
	return &Helper{
		Name: name,
		Fn:   fn,
	}
}

// Execute 执行
func Execute(input []byte, helpers []*Helper) (output []byte, err error) {
	if len(helpers) == 0 {
		return input, nil
	}
	for _, helper := range helpers {
		output, err = helper.Fn(input)
		if err != nil {
			return nil, fmt.Errorf("helper=%q has error:%w", helper.Name, err)
		}
	}
	return output, err
}

// Defaults 默认的helper方法
var Defaults = []*Helper{
	New("osenv", OsEnvVars),
}
