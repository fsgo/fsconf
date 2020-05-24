/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"testing"

	"github.com/fsgo/fsenv"

	"github.com/fsgo/fsconf/internal/helper"
	"github.com/fsgo/fsconf/internal/parser"
)

func Test_confImpl(t *testing.T) {
	conf := New()
	env := fsenv.NewAppEnv(fsenv.Value{
		RootDir: "./testdata",
	})
	conf.SetEnvOnce(env)
	var a interface{}
	if err := conf.Parse("abc.json", &a); err == nil {
		t.Errorf("expect has error")
	}

	conf.RegisterParser(parser.FileJSON, parser.JSON)

	if err := conf.Parse("abc.xyz", &a); err == nil {
		t.Errorf("expect has error 2")
	}
}

func TestNewDefault1(t *testing.T) {
	hd := append([]*helper.Helper{}, helper.Defaults...)
	defer func() {
		helper.Defaults = hd
		if re := recover(); re == nil {
			t.Errorf("want panic")
		}
	}()
	h := helper.New("test", helper.OsEnvVars)
	// helper 有重复的时候
	helper.Defaults = append(helper.Defaults, h, h)
	NewDefault()
}
