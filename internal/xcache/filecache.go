// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/11

package xcache

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileCache struct {
	Dir string
}

func (fc *FileCache) filePath(key string) string {
	m5 := md5.New()
	m5.Write([]byte(key))
	k := fmt.Sprintf("%x", m5.Sum(nil))
	return filepath.Join(fc.Dir, k)
}

func (fc *FileCache) Get(key string, ttl time.Duration) ([]byte, bool) {
	fp := fc.filePath(key)
	bf, err := os.ReadFile(fp)
	if err != nil {
		return nil, false
	}
	if ttl <= 0 {
		return bf, true
	}
	info, err := os.Stat(fp)
	if err != nil {
		return nil, false
	}
	if info.ModTime().Add(ttl).Before(time.Now()) {
		return nil, false
	}
	return bf, true
}

func (fc *FileCache) Set(key string, value []byte) {
	fp := fc.filePath(key)
	_ = os.Remove(fp)
	if err := os.WriteFile(fp, value, 0644); err == nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(fp), os.ModePerm)
	_ = os.WriteFile(fp, value, 0644)
}
