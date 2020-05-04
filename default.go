/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

// Default 默认的实例
var Default IConf = NewDefault()

// Parse 解析配置，配置文件默认认为在 conf/目录下,
// 如 有 conf/abc.toml ，则 confName="abc.toml"
func Parse(confName string, obj interface{}) (err error) {
	return Default.Parse(confName, obj)
}

// ParseByAbsPath 解析绝对路径的配置
func ParseByAbsPath(confAbsPath string, obj interface{}) (err error) {
	return Default.ParseByAbsPath(confAbsPath, obj)
}

// Exists  判断是否存在
func Exists(confName string) bool {
	return Default.Exists(confName)
}

// RegisterParser 注册一个解析器
func RegisterParser(fileExt string, fn ParserFn) error {
	return Default.RegisterParser(fileExt, fn)
}

// RegisterHelper 注册一个辅助方法
func RegisterHelper(h *Helper) error {
	return Default.RegisterHelper(h)
}
