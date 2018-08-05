package main

import (
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/pkg/errors"
)

type StoreManager struct {
	connected bool
	Engine    *xorm.Engine
}

func NewStoreManager() *StoreManager {
	return &StoreManager{connected:false}
}

func (sm *StoreManager) Init() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.PqUser,
		config.PqPwd, config.PqDomain, config.PqDb)
	var err error = nil
	sm.Engine, err = xorm.NewEngine("postgres", connStr)

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




