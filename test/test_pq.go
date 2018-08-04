package main

import (
	_ "github.com/lib/pq"
	"github.com/go-xorm/xorm"
	"fmt"
)

type MessageBody struct {
	MsgType  int
	From     string
	To       string
	Time     int64
	MimeType string
	Content  string
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "postgres",
		"Admin@123", "localhost", "leoim")
	e, err := xorm.NewEngine("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	err = e.Ping()
	if err != nil {
		fmt.Println(err)
	}
	m := new(MessageBody)

	exist, err := e.IsTableExist(m)
	if !exist {
		err = e.CreateTables(m)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(exist)
}