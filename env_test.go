/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/2
 */

package fsconf

import (
	"testing"
)

func TestConfEnv_SetConfRootPath1(t *testing.T) {
	env := &Env{}
	want := "./testdata/"
	env.SetConfRootPath(want)
	got := env.ConfRootPath()
	if want != got {
		t.Errorf("got=%q,want=%q", got, want)
	}
}
func TestConfEnv_SetConfRootPath2(t *testing.T) {
	defer func() {
		if re := recover(); re == nil {
			t.Errorf("want panic")
		}
	}()
	env := &Env{}
	want := "./testdata/"
	env.SetConfRootPath(want)

	env.SetConfRootPath("conf")
}
