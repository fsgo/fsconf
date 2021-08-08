// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/8

package fsconf

import (
	"context"
	"fmt"
	"regexp"

	"github.com/fsgo/fsenv"
)

var _ Hook = (*fsEnvHelper)(nil)

type fsEnvHelper struct {
}

func (f *fsEnvHelper) Name() string {
	return "fsenv"
}

// 模板变量格式：{fsenv.变量名}
var fsEnvVarReg = regexp.MustCompile(`\{fsenv\.([A-Za-z0-9_]+)\}`)

func (f *fsEnvHelper) Execute(ctx context.Context, p *HookParam) (output []byte, err error) {
	contentNew := fsEnvVarReg.ReplaceAllFunc(p.Content, func(subStr []byte) []byte {
		// 将 {fsenv.xxx} 中的 xxx 部分取出
		key := subStr[len("{fsenv.") : len(subStr)-1] // eg: xxx
		var val string
		val, err = f.getValue(string(key), p.Configure)
		if err != nil {
			return nil
		}
		return []byte(val)
	})
	if err != nil {
		return nil, err
	}
	return contentNew, err
}

func (f *fsEnvHelper) getValue(key string, cf Configure) (string, error) {
	cae, ok := cf.(fsenv.HasAppEnv)
	if !ok {
		return "", fmt.Errorf("cannot get appenv")
	}
	ae := cae.AppEnv()

	var value string
	switch key {
	case "RootDir":
		value = ae.RootDir()
	case "IDC":
		value = ae.IDC()
	case "DataRootDir":
		value = ae.DataRootDir()
	case "ConfRootDir":
		value = ae.ConfRootDir()
	case "LogRootDir":
		value = ae.LogRootDir()
	case "RunMode":
		value = string(ae.RunMode())
	default:
		return "", fmt.Errorf("key=%q not support", key)
	}
	return value, nil
}
