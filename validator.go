// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/8/13

package fsconf

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validator 自动规则校验器
//
// 默认使用的 github.com/go-playground/validator/v10
// 如下设置所有字段都是必填的：
//
//	type Address struct {
//		Street string `validate:"required"`
//		City   string `validate:"required"`
//		Planet string `validate:"required"`
//		Phone  string `validate:"required"`
//	}
type Validator interface {
	Validate(val any) error
}

var _ Validator = (*validatorV10)(nil)

var vv10 = newValidatorV10()

func newValidatorV10() *validatorV10 {
	return &validatorV10{
		vv: validator.New(),
	}
}

type validatorV10 struct {
	vv *validator.Validate
}

func (v *validatorV10) Validate(val any) error {
	if v.vv == nil {
		return nil
	}
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Struct:
		return v.vv.Struct(val)
	case reflect.Ptr:
		rvv := rv.Elem()
		switch rvv.Kind() {
		case reflect.Ptr:
			return v.vv.Struct(rvv.Interface())
		case reflect.Struct:
			return v.vv.Struct(val)
		}
	}
	return nil
}

// SetupValidatorV10  配置默认的 validator
//
// 若是有扩展 validator，可以通过这个方法进行设置
func SetupValidatorV10(fn func(v10 *validator.Validate)) {
	fn(vv10.vv)
}
