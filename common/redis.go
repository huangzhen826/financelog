package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisp *redis.Pool

const maxActive = 20

func InitRedis(redisaddr string) {

	//mycfg := jconfig.NewJConfig("redis.cfg")

	if redisp == nil {
		redisp = &redis.Pool{
			MaxIdle:     5,
			MaxActive:   maxActive,
			IdleTimeout: 120 * time.Second,
			Dial: func() (redis.Conn, error) {
				//c, err := redis.Dial("tcp", fmt.Sprint(mycfg.Cfg["ipaddr"]))
				c, err := redis.Dial("tcp", redisaddr)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
	}
}

func GetRedisMaxActive() int {
	return maxActive
}

func RedisZAdd(key, value, score string) bool {
	r := redisp.Get()
	defer r.Close()

	if _, err := r.Do("ZADD", key, score, value); err != nil {
		fmt.Println("redis zadd failed, key:", key, ", err:", err)
		return false
	}
	return true
}

func RedisZRangeByScore(key, score string) ([]string, bool) {
	r := redisp.Get()
	defer r.Close()

	if v, err := redis.Strings(r.Do("ZRANGEBYSCORE", key, 0, score)); err != nil {
		fmt.Println("Redis ZRangeByScore failed, key:", key, ", err:", err)
		return nil, false
	} else {
		return v, true
	}
}

func RedisSet(key, value string, cleanTime int) bool {
	r := redisp.Get()
	defer r.Close()
	if _, err := r.Do("SET", key, value); err != nil {
		fmt.Println("redis set failed, key:", key, ", err:", err)
		return false
	}
	if cleanTime > 0 {
		if _, err := r.Do("EXPIRE", key, fmt.Sprint(cleanTime)); err != nil {
			fmt.Println("redis set expire failed, key:", key)
		}
	}
	return true
}
func RedisGet(key string) (string, bool) {
	r := redisp.Get()
	defer r.Close()
	if v, err := redis.String(r.Do("GET", key)); err != nil {
		fmt.Println("redis get failed, key:", key, ", err:", err)
		return "", false
	} else {
		return v, true
	}
}
func RedisDel(key string) {
	r := redisp.Get()
	defer r.Close()
	r.Do("DEL", key)
}

func GetHashStringMap(k string) (map[string]string, bool) {
	r := redisp.Get()
	defer r.Close()
	v, err := redis.StringMap(r.Do("HGETALL", k))

	if err != nil {
		fmt.Println(err)
		return nil, false
	} else {
		return v, true
	}
}

func SetHashStringMap(k string, m map[string]string) bool {
	r := redisp.Get()
	defer r.Close()
	if _, err := r.Do("HMSET", redis.Args{}.Add(k).AddFlat(m)...); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func RedisHGet(key, childkey string) (string, bool) {
	r := redisp.Get()
	defer r.Close()
	if v, err := redis.String(r.Do("HGET", key, childkey)); err != nil {
		fmt.Println("redis hget failed, key:", key, " childkey:", childkey, ", err:", err)
		return "", false
	} else {
		return v, true
	}
}

func RedisHSet(table, key, value string) bool {
	r := redisp.Get()
	defer r.Close()
	if _, err := r.Do("HSET", table, key, value); err != nil {
		fmt.Println("redis hset failed, table:", table, " key:", key, " value:", value, ", err:", err)
		return false
	} else {
		return true
	}
}

func RedisZRem(key, value string) {
	r := redisp.Get()
	defer r.Close()
	r.Do("ZREM", key, value)
}

func Get_md5_string(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func Get_sql_key(sql string) string {
	sql_key := "data_"
	sql_key += Get_md5_string(sql)
	return sql_key
}
