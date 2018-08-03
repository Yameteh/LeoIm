package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return
	}

	reply, err := redis.Strings(c.Do("HMGET", "ggg", "password","token"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(reply))
	fmt.Println(reply)



	defer c.Close()

}
