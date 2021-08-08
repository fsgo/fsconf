# fsconf
## 1.功能概述
一个可扩展的、简单的配置读取库，目前支持`.json`、`.toml`文件的配置。  


[![Build Status](https://travis-ci.org/fsgo/fsconf.png?branch=master)](https://travis-ci.org/fsgo/fsconf)
[![GoCover](https://gocover.io/_badge/github.com/fsgo/fsconf)](https://gocover.io/github.com/fsgo/fsconf)
[![GoDoc](https://godoc.org/github.com/fsgo/fsconf?status.svg)](https://godoc.org/github.com/fsgo/fsconf)


## 2.对外接口
```go
// 读取并解析配置文件
// confName ：相对于 conf/ 目录的文件路径
// 也支持使用绝对路径
Parse(confName string, obj interface{})error

// 使用绝对/相对 读取并解析配置文件
ParseByAbsPath(confAbsPath string, obj interface{})  error

// ParseBytes 解析bytes
// fileExt 是文件后缀，如.json、.toml
ParseBytes(fileExt string,content []byte,obj interface{})error

// 配置文件是否存在
Exists(confName string) bool

// 注册一个指定后缀的配置的parser
// 如要添加 .ini 文件的支持，可在此注册对应的解析函数即可
RegisterParser(fileExt string, fn ParserFn) error

// 注册一个 Hook
RegisterHook(h Helper) error
```

```go
// NewDefault 创建一个新的配置解析实例
// 会注册默认的配置解析方法和辅助方法
func NewDefault() Configure 
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

###  4.1 hook:从系统环境变量读取变量
配置内容：
```toml
# 若环境变量里有 server_port，而且不为空，则使用环境变量的值，否则使用默认值8080
port = "{osenv.server_port|8080}"

port2 = "{osenv.server_port2}"
```
这样就可以在运行前通过设置环境变量来影响配置文件：
```
export  server_port=80
go run main.go
```

### 4.2 设置配置读取路径
考虑到不同子模块读取配置的目录可能不同，允许让模块自己设置读取配置文件的根目录。
```go
conf:=fsconf.NewDefault()
env:=fsenv.NewAppEnv(fsenv.Value{RootDir:"./testdata/"})
conf.SetEnv(env)
// your code
var confData map[string]string
conf.Parse("abc.json",&confData)
```

### 4.3 .json格式配置
配置注释：每行以`#`开头的是注释，在解析时会忽略掉，如：
```javascript
{
    "ID": 1
#这是注释
   # 这也是注释
}
```


###  4.4 hook:从 appenv 读取变量
```toml
# 补充上 app 的log 目录的路径
LogFilePath = "{fsenv.LogRootDir}/http/access.log"
```

支持：
{fsenv.RootDir}、{fsenv.IDC}、{fsenv.DataRootDir}、
{fsenv.ConfRootDir}、{fsenv.LogRootDir}、{fsenv.RunMode} 。
不支持其他的 key，否则将报错