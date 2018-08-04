package main

import (
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type StoreManager struct {
	engine *xorm.Engine
}

func NewStoreManager() *StoreManager {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.PqUser,
		config.PqPwd, config.PqDomain, config.PqDb)
	e, err := xorm.NewEngine("postgres", connStr)
	if err != nil {
		glog.Error(err)
	}
	return &StoreManager{e}
}

func (sm *StoreManager) Init() error {
	if sm.engine != nil {
		err := sm.engine.Ping()
		if err == nil {
			mb := new(MessageBody)
			exist, err := sm.engine.IsTableExist(mb)
			if !exist {
				sm.engine.CreateTables(mb)
			}
			return err
		} else {
			return err
		}
	} else {
		return errors.New("xorm engine not created")
	}
}

func (sm *StoreManager) Insert(data interface{}) error {
	err := sm.engine.Ping()
	if err == nil {
		if m, ok := data.(*MessageBody);  ok{
			_, err := sm.engine.Insert(m)
			return err
		}
		return errors.New("store insert not support type")
	}else {
		return err
	}
}




