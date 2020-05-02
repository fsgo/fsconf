/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

// IEnv 环境信息
type IEnv interface {
	// 设置应用 配置文件根目录，可选
	SetConfRootPath(confRootPath string)

	// 获取配置文件更目录路径，默认:conf/
	ConfRootPath() string
}

// Env 配置的环境
type Env struct {
	ConfRootDir string
}

// SetConfRootPath 设置配置根目录
func (c *Env) SetConfRootPath(confRootDir string) {
	if c.ConfRootDir != "" {
		panic("cannot set Env.ConfRootDir twice")
	}
	c.ConfRootDir = confRootDir
}

// ConfRootPath  获取配置根目录
func (c *Env) ConfRootPath() string {
	return c.ConfRootDir
}

var _ IEnv = (*Env)(nil)
