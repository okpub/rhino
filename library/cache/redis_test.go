package cache

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func init() {
	config := NewRedis(OptionHost("192.168.0.100"))

	conn := config.Copy(OptionHost("localhost")).Open().Get()
	fmt.Println(redis.Int(conn.Do("get", "uid")))

	conn2 := config.Open().Get()
	fmt.Println(conn2.Err())

	fmt.Println(redis.Int(conn2.Do("get", "uid")))
	fmt.Println("over")
}
