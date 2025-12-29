package config

import (
	"github.com/Verthandii/spring/config/source"
	"github.com/Verthandii/spring/encoding"
	"github.com/Verthandii/spring/ilogger"
)

// 默认为json配置文件
var codec = encoding.GetCodec("json")

var _cfg *Config

type Config struct {
	src   source.Source
	iconf any
}

func SetDefaultCodec(c encoding.Codec) {
	codec = c
}

// Get 获取配置
func Get() any {
	return _cfg.iconf
}

// Init 初始化配置
func Init(src source.Source, iconf any) error {
	_cfg = &Config{
		src:   src,
		iconf: iconf,
	}
	err := _cfg.load()
	if err != nil {
		return err
	}
	ilogger.Info("【Config】初始化成功", "Data", _cfg.iconf)
	return nil
}

// load 加载配置文件到内存中
func (c *Config) load() error {
	sd, err := c.src.Load()
	if err != nil {
		return err
	}
	if err = c.setupConfig(sd); err != nil {
		return err
	}
	return nil
}

// setupConfig 反序列化配置数据
func (c *Config) setupConfig(data []byte) error {
	return codec.Unmarshal(data, c.iconf)
}
