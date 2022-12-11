// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/11

package xcache

import (
	"net/url"
	"time"
)

type Param struct {
	// Timeout 获取数据的超时时间
	Timeout time.Duration

	// TTL 获取数据失败后，若缓存数据在此有效期，则使用缓存
	// > 0 是缓存有效
	TTL time.Duration
}

func ParserParam(str string) (*Param, error) {
	p := &Param{}

	values, err1 := url.ParseQuery(str)
	if err1 != nil {
		return nil, err1
	}

	var err error
	p.Timeout, err = getTimeValue(values, "timeout")

	if err != nil {
		return nil, err
	}

	p.TTL, err = getTimeValue(values, "cache")
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getTimeValue(vs url.Values, key string) (time.Duration, error) {
	v := vs.Get(key)
	if len(v) == 0 {
		return 0, nil
	}
	return time.ParseDuration(v)
}
