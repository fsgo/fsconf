/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/25
 */

package fsconf_test

import (
	"fmt"
	"log"

	"github.com/fsgo/fsconf"
)

func ExampleParseBytes() {
	type User struct {
		Name string
		Age  int
	}
	content := []byte(`{"Name":"Hello","age":18}`)

	var user *User
	if err := fsconf.ParseBytes(".json", content, &user); err != nil {
		log.Fatalln("ParseBytes with error:", err)
	}
	fmt.Println("Name=", user.Name)
	fmt.Println("Age=", user.Age)
	// OutPut:
	// Name= Hello
	// Age= 18
}
