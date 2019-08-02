package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func init() {
	config := NewRedis()

	t := time.Now()
	conn := config.Copy(OptionHost("localhost")).Open().Get()
	fmt.Println(redis.Int(conn.Do("get", "uid")))

	conn2 := config.Open().Get()
	fmt.Println("is error:", conn2.Err())

	fmt.Println(redis.Int(conn2.Do("get", "uid")))
	fmt.Println("over", time.Since(t))
}
