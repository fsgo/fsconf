# fsconf
## 1.功能概述
基于json、toml的，一个简单的配置读取库。

## 2.对外接口：
```go
// 读取并解析配置文件
// confName 不包括 conf/ 目录的文件路径
Parse(confName string, obj interface{}) (err error)

// 使用绝对/相对 读取并解析配置文件
ParseByAbsPath(confAbsPath string, obj interface{}) (err error)

// 配置文件是否存在
Exists(confName string) bool

// 注册一个指定后缀的配置的parser
// 如要添加 .ini 文件的支持，可在此注册对应的解析函数即可
RegisterParser(fileExt string, fn ParserFn) error

// 注册一个辅助方法
RegisterHelper(name string, fn HelperFn) error
```

```go
// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault() IConf 
```

## 3.使用示例

```go
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
    // 默认是从 conf 目录里读取配置
	if err := fsconf.Parse("hosts.json", &hs); err != nil {
		log.Fatal(err)
	}

	fmt.Println("hosts:", hs)
}

```

## 4.特性说明

###  4.1 从系统环境变量读取变量
配置内容：
```
# 若环境变量里有 server_port，而且不为空，则使用环境变量的值，否则使用默认值8080
port = "{osenv.server_port|8080}"

port2 = "{osenv.server_port2}"
```

### 4.2 设置配置读取路径
考虑到不同子模块读取配置的目录可能不同，允许让模块自己设置读取配置文件的根目录。
```go
conf:=fsconf.NewDefault()
env:=fsenv.NewAppEnv(&fsenv.Value{RootDir:"./testdata/"})
conf.SetEnvOnce(env)
// your code
var confData map[string]string
conf.Parse("abc.json",&confData)
```