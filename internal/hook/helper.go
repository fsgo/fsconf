// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/4

package hook

// Fn helper 的函数
type Fn func(cfPath string, confContent []byte) ([]byte, error)
