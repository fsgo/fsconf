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

var _ Hook = (*hookFsEnv)(nil)

type hookFsEnv struct{}

func (f *hookFsEnv) Name() string {
	return "fsenv"
}

// 模板变量格式：{fsenv.变量名}
var fsEnvVarReg = regexp.MustCompile(`\{fsenv\.([A-Za-z0-9_]+)\}`)

func (f *hookFsEnv) Execute(ctx context.Context, p *HookParam) (output []byte, err error) {
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

func (f *hookFsEnv) getValue(key string, _ *Configure) (string, error) {
	var value string
	switch key {
	case "RootDir":
		value = fsenv.RootDir()
	case "IDC":
		value = fsenv.IDC()
	case "DataDir":
		value = fsenv.DataDir()
	case "ConfDir":
		value = fsenv.ConfDir()
	case "TempDir":
		value = fsenv.TempDir()
	case "LogDir":
		value = fsenv.LogDir()
	case "RunMode":
		value = fsenv.RunMode().String()
	default:
		return "", fmt.Errorf("key=%q not support", key)
	}
	return value, nil
}
