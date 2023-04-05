// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/4/5

package fsconf

func inSlice[T comparable](vs []T, val T) bool {
	for i := 0; i < len(vs); i++ {
		if vs[i] == val {
			return true
		}
	}
	return false
}
