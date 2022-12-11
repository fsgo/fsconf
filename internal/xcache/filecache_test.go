// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/11

package xcache

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFileCache_Get(t *testing.T) {
	fc := &FileCache{
		Dir: filepath.Join("testdata", "cache"),
	}
	val1 := []byte("hello")
	fc.Set("k1", val1)
	got1, ok1 := fc.Get("k1", 0)
	require.True(t, ok1)
	require.Equal(t, string(val1), string(got1))

	got1, ok1 = fc.Get("k1", time.Hour)
	require.True(t, ok1)
	require.Equal(t, string(val1), string(got1))
}
