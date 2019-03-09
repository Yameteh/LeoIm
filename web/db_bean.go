package main

type Admin struct {
	Account  string
	Password string
	Level    int8
	Nick     string
}

type OnlineUser struct {
	Id        int64  `xorm:"pk autoincr"`
	Account   string `xorm:"unique"`
	Domain    string
	LoginTime int64
	Level     int8
}
