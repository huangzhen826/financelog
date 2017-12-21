package common

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisNode struct {
	pool redis.Pool
}

var delScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end`)

func (rn *RedisNode) Set(key, value string, expiry time.Duration) bool {
	conn := rn.pool.Get()
	defer conn.Close()

	var result string

	if expiry > 0 {
		result, _ = redis.String(conn.Do("SET", key, value, "NX", "EX", int(expiry/time.Millisecond)))
	} else {
		result, _ = redis.String(conn.Do("SET", key, value, "NX"))
	}

	if result != "OK" {
		return false
	}
	return true

}

func (rn *RedisNode) Del(key, value string) bool {
	conn := rn.pool.Get()
	defer conn.Close()

	delScript.Do(conn, key, value)

	return true
}

func NewRedisMuxtex() *RedisNode {
	return &RedisNode{
		pool: *redisp,
	}
}
