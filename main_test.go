// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/9/17

package fsconf

import (
	"os"

	"github.com/fsgo/fsenv"
)

func init() {
	_ = os.Setenv("Port1", "8080")
	_ = os.Setenv("Port2", "8081")
	_ = os.Setenv("APP", "demo.fenji")
	env := fsenv.NewAppEnv(fsenv.Value{RootDir: "./testdata"})
	Default.(fsenv.CanSetAppEnv).SetAppEnv(env)
}
