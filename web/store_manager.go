package main

import (
	"fmt"

	"github.com/golang/glog"

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

	fmt.Println(connStr)
	var err error = nil
	sm.Engine, err = xorm.NewEngine("mysql", connStr)

	if err == nil {
		err = sm.Engine.Ping()
		if err == nil {
			sm.connected = true
			sm.setupTables()
		}
	}
	return err
}

func (sm *StoreManager) setupTables() {
	if sm.connected {
		a := new(Admin)
		exist, err := sm.Engine.IsTableExist(a)
		if !exist {
			sm.Engine.CreateTables(a)
		} else {
			glog.Error(err)
		}
	}
}

func (sm *StoreManager) CheckAdmin(a string, p string) bool {
	admin := &Admin{Account: a, Password: p}
	has, err := sm.Engine.Get(admin)
	if err != nil {
		glog.Error(err)
		return false
	}
	return has
}
