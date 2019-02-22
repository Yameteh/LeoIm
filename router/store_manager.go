package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

type StoreManager struct {
	connected bool
	Engine    *xorm.Engine
}

func NewStoreManager() *StoreManager {
	return &StoreManager{connected: false}
}

func (sm *StoreManager) Init() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?sslmode=disable", config.MysqlUser,
		config.MysqlPwd, config.Domain, config.MysqlPort, config.MysqlDb)
	var err error = nil
	sm.Engine, err = xorm.NewEngine("mysql", connStr)

	if err == nil {
		err = sm.Engine.Ping()
		if err == nil {
			sm.connected = true
			mb := new(MessageBody)
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
			if m, ok := data.(*MessageBody); ok {
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

func (sm *StoreManager) QueryMessage(user string, time int64, out *[]MessageBody) error {
	if sm.Engine != nil {
		if sm.connected {
			err := sm.Engine.Where("message_body.to=? AND message_body.time>=?", user, time).Find(out)
			return err
		} else {
			return errors.New("xorm query message when not connect database")
		}
	} else {
		return errors.New("xorm query message when engine is nil")
	}
}
