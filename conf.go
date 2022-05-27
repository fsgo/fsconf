// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsgo/fsenv"
)

// Configure 配置解析定义
type Configure interface {
	// Parse 读取并解析配置文件
	// confName ：相对于 conf/ 目录的文件路径
	// 也支持使用绝对路径
	Parse(confName string, obj interface{}) error

	// ParseByAbsPath 使用绝对/相对 读取并解析配置文件
	ParseByAbsPath(confAbsPath string, obj interface{}) error

	// ParseBytes 解析bytes内容
	ParseBytes(fileExt string, content []byte, obj interface{}) error

	// Exists 配置文件是否存在
	Exists(confName string) bool

	// RegisterParser 注册一个指定后缀的配置的 parser
	// 如要添加 .ini 文件的支持，可在此注册对应的解析函数即可
	RegisterParser(fileExt string, fn ParserFn) error

	// RegisterHook 注册一个辅助方法
	RegisterHook(h Hook) error

	// WithContext 设置一个 context，并返回新的对象
	WithContext(ctx context.Context) Configure
}

// New 创建一个新的配置解析实例
// 返回的实例是没有注册任何解析能力的
func New() Configure {
	conf := &confImpl{
		parsers: map[string]ParserFn{},
	}
	return conf
}

// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault() Configure {
	conf := New()
	for name, fn := range defaultParsers {
		if err := conf.RegisterParser(name, fn); err != nil {
			panic(fmt.Sprintf("RegisterParser(%q) err=%s", name, err))
		}
	}

	for _, h := range defaultHooks {
		if err := conf.RegisterHook(h); err != nil {
			panic(fmt.Sprintf("RegisterInterceptor(%q) err=%s", h.Name(), err))
		}
	}
	return conf
}

type confImpl struct {
	fsenv.WithAppEnv
	parsers map[string]ParserFn
	hooks   hooks
	ctx     context.Context
}

func (c *confImpl) Parse(confName string, obj interface{}) (err error) {
	confAbsPath, err := c.confFileAbsPath(confName)
	if err != nil {
		return err
	}
	return c.ParseByAbsPath(confAbsPath, obj)
}

func (c *confImpl) confFileAbsPath(confName string) (string, error) {
	if strings.HasPrefix(confName, "./") {
		return filepath.Abs(confName)
	}
	if strings.HasPrefix(confName, "../") {
		return filepath.Abs(confName)
	}
	if filepath.IsAbs(confName) {
		return confName, nil
	}
	return filepath.Join(c.AppEnv().ConfRootDir(), confName), nil
}

func (c *confImpl) ParseByAbsPath(confAbsPath string, obj interface{}) (err error) {
	if len(c.parsers) == 0 {
		return fmt.Errorf("no parser")
	}

	return c.readConfDirect(confAbsPath, obj)
}

func (c *confImpl) readConfDirect(confPath string, obj interface{}) error {
	content, errIO := ioutil.ReadFile(confPath)
	if errIO != nil {
		return errIO
	}
	fileExt := filepath.Ext(confPath)
	return c.parseBytes(confPath, fileExt, content, obj)
}

func (c *confImpl) context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *confImpl) ParseBytes(fileExt string, content []byte, obj interface{}) error {
	return c.parseBytes("", fileExt, content, obj)
}

func (c *confImpl) parseBytes(confPath string, fileExt string, content []byte, obj interface{}) error {
	parserFn, hasParser := c.parsers[fileExt]
	if fileExt == "" || !hasParser {
		return fmt.Errorf("fileExt %q is not supported yet", fileExt)
	}

	p := &HookParam{
		Configure: c,
		ConfPath:  confPath,
		Content:   content,
	}

	contentNew, errHook := c.hooks.Execute(c.context(), p)

	if errHook != nil {
		return errHook
	}

	if errParser := parserFn(contentNew, obj); errParser != nil {
		return fmt.Errorf("%w, config content=\n%s", errParser, string(contentNew))
	}
	return nil
}

func (c *confImpl) Exists(confName string) bool {
	p, err := c.confFileAbsPath(confName)
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
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

func (c *confImpl) RegisterHook(h Hook) error {
	return c.hooks.Add(h)
}

func (c *confImpl) clone() *confImpl {
	c1 := &confImpl{
		parsers: make(map[string]ParserFn, len(c.parsers)),
	}
	for n, fn := range c.parsers {
		c1.parsers[n] = fn
	}
	c1.hooks = append([]Hook{}, c.hooks...)

	if env := c.AppEnv(); env != fsenv.Default {
		c1.SetAppEnv(env)
	}
	return c1
}

func (c *confImpl) WithContext(ctx context.Context) Configure {
	c1 := c.clone()
	c1.ctx = ctx
	return c1
}

var _ Configure = (*confImpl)(nil)
