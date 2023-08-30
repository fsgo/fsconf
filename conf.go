// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package fsconf

import (
	"context"
	"errors"
	"fmt"
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
	Parse(confName string, obj any) error

	// ParseByAbsPath 使用绝对/相对 读取并解析配置文件
	ParseByAbsPath(confAbsPath string, obj any) error

	// ParseBytes 解析bytes内容
	ParseBytes(fileExt string, content []byte, obj any) error

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

// AutoChecker 当配置解析完成后，用于自动校验，
// 这个方法是在 validator 校验完成之后才执行的
type AutoChecker interface {
	AutoCheck() error
}

// New 创建一个新的配置解析实例
// 返回的实例是没有注册任何解析能力的
func New() Configure {
	conf := &confImpl{
		parsers:  map[string]ParserFn{},
		validate: vv10,
	}
	return conf
}

// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault() Configure {
	conf := New()
	for _, pair := range defaultParsers {
		if err := conf.RegisterParser(pair.Name, pair.Fn); err != nil {
			panic(fmt.Sprintf("RegisterParser(%q) err=%s", pair.Name, err))
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
	ctx        context.Context
	validate   Validator
	parsers    map[string]ParserFn
	parseNames []string // 支持的文件后缀，如 []string{".json",".toml"}
	hooks      hooks
	fsenv.WithAppEnv
}

func (c *confImpl) Parse(confName string, obj any) (err error) {
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

func (c *confImpl) ParseByAbsPath(confAbsPath string, obj any) (err error) {
	if len(c.parsers) == 0 {
		return errors.New("no parser")
	}

	return c.readConfDirect(confAbsPath, obj)
}

func (c *confImpl) realConfPath(confPath string) (path string, ext string, err error) {
	fileExt := filepath.Ext(confPath)
	info, err1 := os.Stat(confPath)

	if err1 == nil && !info.IsDir() {
		return confPath, fileExt, nil
	}

	notExist := err1 != nil && os.IsNotExist(err1)
	isDir := err1 == nil && info.IsDir()

	// fileExt == "" 是为了兼容存在同名目录的情况
	if (notExist || isDir || fileExt == "") && !inSlice(c.parseNames, fileExt) {
		for i := 0; i < len(c.parseNames); i++ {
			ext2 := c.parseNames[i]
			name2 := confPath + ext2
			info2, err2 := os.Stat(name2)
			if err2 == nil && !info2.IsDir() {
				return name2, ext2, nil
			}
		}
	}
	if err1 != nil {
		return "", "", err1
	}
	return "", "", fmt.Errorf("cannot get real path for %q", confPath)
}

func (c *confImpl) readConfDirect(confPath string, obj any) error {
	realFile, fileExt, err := c.realConfPath(confPath)
	if err != nil {
		return err
	}
	content, errIO := os.ReadFile(realFile)
	if errIO != nil {
		return errIO
	}
	err2 := c.parseBytes(realFile, fileExt, content, obj)
	if err2 == nil {
		return nil
	}
	return fmt.Errorf("parser %q failed: %w", realFile, err2)
}

func (c *confImpl) context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *confImpl) ParseBytes(fileExt string, content []byte, obj any) error {
	return c.parseBytes("", fileExt, content, obj)
}

func (c *confImpl) parseBytes(confPath string, fileExt string, content []byte, obj any) error {
	parserFn, hasParser := c.parsers[fileExt]
	if len(fileExt) == 0 || !hasParser {
		err1 := fmt.Errorf("fileExt %q is not supported yet", fileExt)
		if confPath == "" {
			return err1
		}
		return fmt.Errorf("cannot parser %q: %w", confPath, err1)
	}

	p := &HookParam{
		FileExt:   fileExt,
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

	if err := c.validate.Validate(obj); err != nil {
		return err
	}

	if ac, ok := obj.(AutoChecker); ok {
		if err := ac.AutoCheck(); err != nil {
			return fmt.Errorf("autoCheck: %w", err)
		}
	}
	return nil
}

func (c *confImpl) Exists(confName string) bool {
	p, err := c.confFileAbsPath(confName)
	if err != nil {
		return false
	}

	info, err := os.Stat(p)
	if err == nil && !info.IsDir() {
		return true
	}
	if !os.IsNotExist(err) {
		return false
	}
	for ext := range c.parsers {
		info1, err1 := os.Stat(p + ext)
		if err1 == nil && !info1.IsDir() {
			return true
		}
	}
	return false
}

func (c *confImpl) RegisterParser(fileExt string, fn ParserFn) error {
	if _, has := c.parsers[fileExt]; has {
		return fmt.Errorf("parser=%q already exists", fileExt)
	}
	c.parsers[fileExt] = fn
	c.parseNames = append(c.parseNames, fileExt)
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
