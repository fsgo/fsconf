// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package helper

import (
	"bytes"
	"os"
	"regexp"
)

// 模板变量格式：{env.变量名} 或者 {env.变量名|默认值}
var osEnvVarReg = regexp.MustCompile(`\{osenv\.([A-Za-z0-9_]+)(\|[^}]+)?\}`)

// OsEnvVars 将配置文件中的 {env.xxx} 的内容，从环境变量中读取并替换
func OsEnvVars(content []byte) ([]byte, error) {
	contentNew := osEnvVarReg.ReplaceAllFunc(content, func(subStr []byte) []byte {
		// 将 {osenv.xxx} 中的 xxx 部分取出
		// 或者 将 {osenv.yyy|val} 中的 yyy|val 部分取出

		keyWithDefaultVal := subStr[len("{osenv.") : len(subStr)-1] // eg: xxx 或者 yyy|val
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
