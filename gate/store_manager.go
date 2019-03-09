package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type StoreManager struct {
	connected bool
	Engine    *xorm.Engine
}

func NewStoreManager() *StoreManager {
	return &StoreManager{connected: false}
}

func (sm *StoreManager) Init() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config.MysqlUser,
		config.MysqlPwd, config.MysqlDomain, config.MysqlPort, config.MysqlDb)
	var err error = nil
	sm.Engine, err = xorm.NewEngine("mysql", connStr)

	if err == nil {
		err = sm.Engine.Ping()
		if err == nil {
			sm.connected = true
			mb := new(OnlineUser)
			var exist bool
			exist, err = sm.Engine.IsTableExist(mb)
			if !exist {
				sm.Engine.CreateTables(mb)
			}
		}
	}
	return err
}

func (sm *StoreManager) Insert(data interface{}) error {
	if sm.Engine != nil {
		if sm.connected {
			if m, ok := data.(*OnlineUser); ok {
				_, err := sm.Engine.Insert(m)
				return err
			}
			return errors.New("xorm engine insert not support data type")
		} else {
			return errors.New("xorm engine insert when not connect database")
		}
	} else {
		return errors.New("xorm insert when engine is nil")
	}
}

func (sm *StoreManager) Delete(data interface{}) error {
	if sm.Engine != nil {
		if sm.connected {
			if m, ok := data.(*OnlineUser); ok {
				_, err := sm.Engine.Delete(m)
				return err
			}
			return errors.New("xorm engine delete not support data type")
		} else {
			return errors.New("xorm engine delete when not connected database")
		}
	} else {
		return errors.New("xorm delete when engine is nil")
	}
}
