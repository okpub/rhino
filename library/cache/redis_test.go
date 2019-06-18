package cache

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func init() {
	pool := NewRedis(OptionAddr("192.168.0.100:6379")).Open()
	conn := pool.Get()
	n, _ := redis.Int(conn.Do("get", "uid"))
	fmt.Println(n)
}
