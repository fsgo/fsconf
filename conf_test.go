/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"testing"
)

func Test_confImpl(t *testing.T) {
	conf := New(&ConfEnv{
		RootPath: "./testdata/",
	})
	var a interface{}
	if err := conf.Parse("abc.json", &a); err == nil {
		t.Errorf("expect has error")
	}

	conf.RegisterParser(FileExtJSON, JSONParser)

	if err := conf.Parse("abc.xyz", &a); err == nil {
		t.Errorf("expect has error 2")
	}
}
