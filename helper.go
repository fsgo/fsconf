/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
)

// HelperFn 辅助功能方法
// 在正式解析配置前执行
type HelperFn func(confContent []byte) ([]byte, error)

type helper struct {
	name string
	fn   HelperFn
}

// 以下为默认的helper 方法

// 模板变量格式：{env.变量名} 或者 {env.变量名|默认值}
var helperOsEnvVarReg = regexp.MustCompile(`\{osenv\.([A-Za-z0-9_]+)(\|[^}]+)?\}`)

// 将配置文件中的 {env.xxx} 的内容，从环境变量中读取并替换
func helperOsEnvVars(content []byte) ([]byte, error) {
	contentNew := helperOsEnvVarReg.ReplaceAllFunc(content, func(subStr []byte) []byte {
		// 将 {osenv.xxx} 中的 xxx 部分取出
		// 或者 将 {osenv.yyy|val} 中的 yyy|val 部分取出

		keyWithDefaultVal := subStr[7 : len(subStr)-1] // eg: xxx 或者 yyy|val
		idx := bytes.Index(keyWithDefaultVal, []byte("|"))
		if idx > 0 {
			// {osenv.变量名|默认值} 有默认值的格式
			key := string(keyWithDefaultVal[:idx])  // eg: yyy
			defaultVal := keyWithDefaultVal[idx+1:] // eg: val
			envVal := os.Getenv(key)
			if envVal == "" {
				return defaultVal
			}
			return []byte(envVal)
		}

		// {osenv.变量名} 无默认值的部分
		return []byte(os.Getenv(string(keyWithDefaultVal)))
	})
	return contentNew, nil
}

func executeHelpers(input []byte, helpers []*helper) (output []byte, err error) {
	if len(helpers) == 0 {
		return input, nil
	}
	for _, helper := range helpers {
		output, err = helper.fn(input)
		if err != nil {
			return nil, fmt.Errorf("helper=%q has error:%w", helper.name, err)
		}
	}
	return output, err
}

var defaultHelpers = map[string]HelperFn{
	"osenv": helperOsEnvVars,
}
