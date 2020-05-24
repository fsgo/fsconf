/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fsgo/fsenv"

	"github.com/fsgo/fsconf/internal/helper"
	"github.com/fsgo/fsconf/internal/parser"
)

// ParserFn 针对特定文件后缀的配置解析方法
// 当前已经内置了 .toml  和 .json的解析方法
type ParserFn = parser.Fn

// Helper 辅助方法，在执行解析前，会先会配置的内容进行解析处理
type Helper = helper.Helper

// IConf 配置解析定义
type IConf interface {
	fsenv.IModuleEnv

	// 读取并解析配置文件
	// confName 不包括 conf/ 目录的文件路径
	Parse(confName string, obj interface{}) (err error)

	// 使用绝对/相对 读取并解析配置文件
	ParseByAbsPath(confAbsPath string, obj interface{}) (err error)

	// 配置文件是否存在
	Exists(confName string) bool

	// 注册一个指定后缀的配置的parser
	// 如要添加 .ini 文件的支持，可在此注册对应的解析函数即可
	RegisterParser(fileExt string, fn ParserFn) error

	// 注册一个辅助方法
	RegisterHelper(h *Helper) error
}

// New 创建一个新的配置解析实例
// 返回的实例是没有注册任何解析能力的
func New() IConf {
	conf := &confImpl{
		parsers:   map[string]ParserFn{},
		ModuleEnv: &fsenv.ModuleEnv{},
	}
	return conf
}

// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault() IConf {
	conf := New()
	for name, fn := range parser.Defaults {
		if err := conf.RegisterParser(name, fn); err != nil {
			panic(fmt.Sprintf("RegisterParser(%q) err=%s", name, err))
		}
	}

	for _, h := range helper.Defaults {
		if err := conf.RegisterHelper(h); err != nil {
			panic(fmt.Sprintf("RegisterHelper(%q) err=%s", h.Name, err))
		}
	}
	return conf
}

type confImpl struct {
	*fsenv.ModuleEnv
	parsers map[string]ParserFn
	helpers []*helper.Helper
}

func (c *confImpl) Parse(confName string, obj interface{}) (err error) {
	confAbsPath := c.confFileAbsPath(confName)
	return c.ParseByAbsPath(confAbsPath, obj)
}

func (c *confImpl) confFileAbsPath(confName string) string {
	return filepath.Join(c.Env().ConfRootDir(), confName)
}

func (c *confImpl) ParseByAbsPath(confAbsPath string, obj interface{}) (err error) {
	if len(c.parsers) == 0 {
		return fmt.Errorf("no parser")
	}

	return c.readConfDirect(confAbsPath, obj)
}

func (c *confImpl) readConfDirect(confPath string, obj interface{}) error {
	fileExt := filepath.Ext(confPath)

	parserFn, hasParser := c.parsers[fileExt]
	if fileExt == "" || !hasParser {
		return fmt.Errorf("fileType '%s' is not yet supported", fileExt)
	}

	content, errIO := ioutil.ReadFile(confPath)
	if errIO != nil {
		return errIO
	}

	contentNew, errHelper := helper.Execute(content, c.helpers)

	if errHelper != nil {
		return errHelper
	}

	if errParser := parserFn(contentNew, obj); errParser != nil {
		return fmt.Errorf("%w, config content=\n%s", errParser, string(contentNew))
	}
	return nil
}

func (c *confImpl) Exists(confName string) bool {
	info, err := os.Stat(c.confFileAbsPath(confName))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (c *confImpl) RegisterParser(fileExt string, fn ParserFn) error {
	if _, has := c.parsers[fileExt]; has {
		return fmt.Errorf("parser=%q already exists", fileExt)
	}
	c.parsers[fileExt] = fn
	return nil
}

func (c *confImpl) RegisterHelper(h *Helper) error {
	if h.Name == "" {
		return fmt.Errorf("name ='', not allow")
	}

	for _, h1 := range c.helpers {
		if h.Name == h1.Name {
			return fmt.Errorf("helper=%q already exists", h.Name)
		}
	}
	c.helpers = append(c.helpers, h)
	return nil
}

var _ IConf = (*confImpl)(nil)
