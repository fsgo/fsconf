// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/9/9

package fsconf_test

import (
	"testing"

	"github.com/fsgo/fst"

	"github.com/fsgo/fsconf"
)

func TestHookTPL_Execute(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		h1 := &fsconf.HookTPL{
			HookName:  "t1",
			KeyPrefix: "t1",
			Values: map[string]string{
				"k1": "v1",
			},
		}
		c1 := fsconf.WithHook(h1)
		content1 := `
{
  "K1" : "{t1.k1}"
}
`
		data1 := map[string]string{}
		err := c1.ParseBytes(".json", []byte(content1), &data1)
		fst.NoError(t, err)
		want1 := map[string]string{"K1": "v1"}
		fst.Equal(t, want1, data1)

		content2 := `
{
  "K1" : "{t1.k2}"
}
`
		data2 := map[string]string{}
		err2 := c1.ParseBytes(".json", []byte(content2), &data2)
		fst.Error(t, err2)
		fst.Empty(t, data2)
	})
}
