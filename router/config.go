package main

import (
	"strings"

	"github.com/golang/glog"
	"gopkg.in/ini.v1"
)

const (
	INI_SECTION_LISTEN = "listen"
	INI_KEY_DOMAIN     = "domain"
	INI_KEY_PORT       = "port"

	INI_SECTION_GATE = "gate_rpc_server"
	INI_KEY_DOMAINS  = "domains"

	INI_SECTION_MYSQL = "mysql_server"
	INI_KEY_DB        = "database"
	INI_KEY_USER      = "user"
	INI_KEY_PASSWORD  = "password"

	INI_SECTION_SYNC_SERVER = "sync_server"

	INI_SECTION_REDIS = "redis_server"
)

type Config struct {
	Domain      string
	Port        int
	GateServer  []string
	MysqlDomain string
	MysqlPort   int
	MysqlDb     string
	MysqlUser   string
	MysqlPwd    string
	WebDomain   string
	WebPort     int
	RedisDomain string
	RedisPort   int
	FileRootDir string
	file        *ini.File
}

func NewConfig(file string) *Config {
	cfg, err := ini.Load(file)
	if err != nil {
		glog.Error("LeoIm router load config ini error : ", err)
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
	}

	s = cfg.file.Section(INI_SECTION_GATE)
	if s != nil {
		if ds, err := s.GetKey(INI_KEY_DOMAINS); err == nil {
			cfg.GateServer = strings.Split(ds.String(), ",")
		} else {
			glog.Error(err)
		}
	}

	s = cfg.file.Section(INI_SECTION_MYSQL)
	if s != nil {
		if d, err := s.GetKey(INI_KEY_DOMAIN); err == nil {
			cfg.MysqlDomain = d.String()
		} else {
			glog.Error(err)
		}

		if p, err := s.GetKey(INI_KEY_PORT); err == nil {
			cfg.MysqlPort, _ = p.Int()
		} else {
			glog.Error(err)
		}

		if db, err := s.GetKey(INI_KEY_DB); err == nil {
			cfg.MysqlDb = db.String()
		} else {
			glog.Error(err)
		}

		if u, err := s.GetKey(INI_KEY_USER); err == nil {
			cfg.MysqlUser = u.String()
		} else {
			glog.Error(err)
		}

		if pwd, err := s.GetKey(INI_KEY_PASSWORD); err == nil {
			cfg.MysqlPwd = pwd.String()
		} else {
			glog.Error(err)
		}

	}

	s = cfg.file.Section(INI_SECTION_SYNC_SERVER)
	if s != nil {
		if d, err := s.GetKey(INI_KEY_DOMAIN); err == nil {
			cfg.WebDomain = d.String()
		} else {
			glog.Error(err)
		}

		if p, err := s.GetKey(INI_KEY_PORT); err == nil {
			cfg.WebPort, _ = p.Int()
		} else {
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
