package confext

import (
	"testing"

	"github.com/fsgo/fsconf"
	"github.com/fsgo/fst"
)

func TestValidator(t *testing.T) {
	type user struct {
		Name string `validate:"required"`
		Age  int
	}

	t.Run("validator-1", func(t *testing.T) {
		var u *user
		err := fsconf.ParseBytes(".json", []byte(`{"Age":12}`), &u)
		fst.Equal(t, &user{Age: 12}, u)
		fst.Error(t, err)
	})

	t.Run("validator-2", func(t *testing.T) {
		var u *user
		err := fsconf.ParseBytes(".json", []byte(``), &u)
		fst.Nil(t, u)
		fst.Error(t, err)
	})

	t.Run("validator-3", func(t *testing.T) {
		var u *user
		err := fsconf.ParseBytes(".json", []byte(`{"Age":12,"Name":""}`), &u)
		fst.Equal(t, &user{Age: 12}, u)
		fst.Error(t, err)
	})

	t.Run("validator-4", func(t *testing.T) {
		var u *user
		err := fsconf.ParseBytes(".json", []byte(`{"Age":12,"Name":"hello"}`), &u)
		fst.Equal(t, &user{Age: 12, Name: "hello"}, u)
		fst.NoError(t, err)
	})

	t.Run("validator-5", func(t *testing.T) {
		u := &user{}
		err := fsconf.ParseBytes(".json", []byte(`{"Age":12,"Name":"hello"}`), u)
		fst.Equal(t, &user{Age: 12, Name: "hello"}, u)
		fst.NoError(t, err)
	})

	t.Run("validator-6", func(t *testing.T) {
		u := &user{}
		err := fsconf.ParseBytes(".json", []byte(`{"Age":12,"Name":""}`), u)
		fst.Equal(t, &user{Age: 12}, u)
		fst.Error(t, err)
	})
}
