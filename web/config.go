package main

import (
	"github.com/golang/glog"
	"gopkg.in/ini.v1"
)

const (
	INI_SECTION_LISTEN     = "listen"
	INI_KEY_DOMAIN         = "domain"
	INI_KEY_PORT           = "port"
	INI_KEY_ROOT_DIRECTORY = "root_directory"
	INI_KEY_SESSION_KEY    = "session_key"
	INI_KEY_FRONT_TITLE    = "front_title"
	INI_SECTION_REDIS      = "redis_server"
	INI_SECTION_MYSQL      = "mysql_server"
	INI_KEY_USER           = "user"
	INI_KEY_PASSWORD       = "password"
	INI_KEY_DATABASE       = "database"
)

type Config struct {
	Domain        string
	Port          int
	RedisDomain   string
	RedisPort     int
	RootDirectory string
	SessionKey    string
	FrontTitle    string
	MysqlDomain   string
	MysqlPort     int
	MysqlUser     string
	MysqlPwd      string
	MysqlDb       string
	file          *ini.File
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

		if rd, err := s.GetKey(INI_KEY_ROOT_DIRECTORY); err == nil {
			cfg.RootDirectory = rd.String()
		} else {
			glog.Error(err)
		}

		if sk, err := s.GetKey(INI_KEY_SESSION_KEY); err == nil {
			cfg.SessionKey = sk.String()
		} else {
			glog.Error(err)
		}

		if ft, err := s.GetKey(INI_KEY_FRONT_TITLE); err == nil {
			cfg.FrontTitle = ft.String()
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

		if db, err := s.GetKey(INI_KEY_DATABASE); err == nil {
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

}
