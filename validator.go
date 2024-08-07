// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/8/13

package fsconf

// AutoChecker 当配置解析完成后，用于自动校验，
// 这个方法是在 validator 校验完成之后才执行的
type AutoChecker interface {
	AutoCheck() error
}

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

var DefaultValidator Validator
