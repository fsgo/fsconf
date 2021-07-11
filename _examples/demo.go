// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu
// Date: 2020/5/2

package main

import (
	"fmt"
	"log"

	"github.com/fsgo/fsconf"
)

type Hosts []Host

type Host struct {
	IP   string
	Port int
}

func main() {
	var hs Hosts
	if err := fsconf.Parse("hosts.json", &hs); err != nil {
		log.Fatal(err)
	}

	fmt.Println("hosts:", hs)
}
