/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

// Env 环境信息
type Env interface {
	// 设置应用 配置文件根目录，可选
	SetConfRootPath(confRootPath string)

	// 获取配置文件更目录路径，默认: RootPath()/conf
	ConfRootPath() string
}

// ConfEnv 配置的环境
type ConfEnv struct {
	RootPath string
}

// SetConfRootPath 设置配置根目录
func (c *ConfEnv) SetConfRootPath(confRootPath string) {
	if c.RootPath != "" {
		panic("cannot set ConfEnv.RootPath twice")
	}
	c.RootPath = confRootPath
}

// ConfRootPath  获取配置根目录
func (c *ConfEnv) ConfRootPath() string {
	return c.RootPath
}

var _ Env = (*ConfEnv)(nil)
