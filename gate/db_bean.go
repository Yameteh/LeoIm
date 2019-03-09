package main

type OnlineUser struct {
	Id        int64  `xorm:"pk autoincr"`
	Account   string `xorm:"unique"`
	Domain    string
	LoginTime int64
	Level     int8
}
