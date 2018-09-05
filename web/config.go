package main

import (
	"gopkg.in/ini.v1"
	"github.com/golang/glog"
)

const (
	INI_SECTION_LISTEN = "listen"
	INI_KEY_DOMAIN = "domain"
	INI_KEY_PORT = "port"
	INI_KEY_ROOT_DIRECTORY = "root_directory"

	INI_SECTION_REDIS = "redis_server"
)

type Config struct {
	Domain     string
	Port       int
	RedisDomain string
	RedisPort  int
	RootDirectory string
	file       *ini.File
}

func NewConfig(file string) *Config {
	cfg, err := ini.Load(file)
	if err != nil {
		glog.Error("LeoIm web load config ini error : ", err)
		return nil
	}
	return &Config{file: cfg}
}

func (cfg *Config) Parse() {
	s := cfg.file.Section(INI_SECTION_LISTEN)
	if s != nil {
		if d, err := s.GetKey(INI_KEY_DOMAIN); err == nil {
			cfg.Domain = d.String()
		} else {
			glog.Error(err)
		}

		if p, err := s.GetKey(INI_KEY_PORT); err == nil {
			cfg.Port, _ = p.Int()
		} else {
			glog.Error(err)
		}

		if rd,err := s.GetKey(INI_KEY_ROOT_DIRECTORY);err == nil {
			cfg.RootDirectory = rd.String()
		}else {
			glog.Error(err)
		}
	}

	s = cfg.file.Section(INI_SECTION_REDIS)
	if s != nil {
		if d, err := s.GetKey(INI_KEY_DOMAIN); err == nil {
			cfg.RedisDomain = d.String()
		} else {
			glog.Error(err)
		}

		if p, err := s.GetKey(INI_KEY_PORT); err == nil {
			cfg.RedisPort, _ = p.Int()
		} else {
			glog.Error(err)
		}

	}

}

