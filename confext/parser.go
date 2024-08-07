package confext

import (
	"github.com/BurntSushi/toml"
	"github.com/fsgo/fsconf"
	"gopkg.in/yaml.v3"
)

type parserNameFn struct {
	Fn   fsconf.DecoderFunc
	Name string
}

// defaultParsers 所有默认的 parser，
// 当传入配置文件名不包含后置的时候，会使用此顺序依次查找
var parsers = []parserNameFn{
	{Name: ".toml", Fn: toml.Unmarshal},
	{Name: ".yml", Fn: yaml.Unmarshal},
	{Name: ".yaml", Fn: yaml.Unmarshal},
}

func init() {
	for _, p := range parsers {
		fsconf.RegisterParser(p.Name, p.Fn)
	}
}
