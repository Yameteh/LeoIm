package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"fmt"
)

type User struct {
	Uuid     string
	Password string
	Token    string
}

func RedisUpdateUser(u *User) error {
	glog.Info("update user ", u.Uuid)
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisDomain, config.RedisPort))
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("HMSET", u.Uuid, "password", u.Password, "token", u.Token)
	return err
}

func RedisQueryUser(uuid string) *User {
	glog.Info("query user ", uuid)
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisDomain, config.RedisPort))
	if err != nil {
		glog.Error(err)
		return nil
	}
	defer c.Close()

	reply, err := redis.Strings(c.Do("HMGET", uuid, "password", "token"))
	if err != nil {
		glog.Error(err)
		return nil
	}

	if len(reply) == 2 {
		user := &User{uuid, reply[0], reply[1]}
		return user
	} else {
		return nil
	}
}

