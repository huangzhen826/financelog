package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func Test_Mutex(t *testing.T) {
	rpool := redis.Pool{
		MaxIdle:     5,
		MaxActive:   20,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "172.16.1.13:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	p := NewRedisMuxtex()

	fmt.Println("pool:", p)

	nodes := []common.Node{
		p,
	}

	m := NewMutex("test", nodes)

	fmt.Println("mutex:", m)

	fmt.Println("main lock:", m.Lock())
	time.Sleep(30000000000)
	fmt.Println("main unlock:", m.Unlock())
}
