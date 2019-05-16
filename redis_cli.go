/**
 * Created by Alen on 2019-05-16 12:18
 */

package models

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"time"
	"github.com/astaxie/beego"
)

const (
	RedisURL            = "redis://127.0.0.1:6379"
	redisMaxIdle        = 3   // 最大空闲连接数
	redisIdleTimeoutSec = 240 // 最大空闲连接时间
	RedisPassword       = ""
)

type redisPool struct {
	pool *redis.Pool // redis connection pool
}

var RedisCli *redisPool // redis client

func init() {
	if RedisCli == nil {
		redisUrl := beego.AppConfig.String("redis_host") + ":" + beego.AppConfig.String("redis_port")
		RedisCli = &redisPool{
			pool: createRedisPool("redis:" + redisUrl),
		}
	}
}

// GetRedisPool返回redis连接池
func createRedisPool(redisURL string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			// 验证redis密码
			//if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
			//	return nil, fmt.Errorf("redis auth password error: %s", authErr)
			//}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func (cli *redisPool) Set(k, v string) {
	c := cli.pool.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func (cli *redisPool) GetStringValue(k string) string {
	c := cli.pool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		fmt.Println("Get Error: ", err.Error())
		return ""
	}
	return username
}

func (cli *redisPool) SetKeyExpire(k string, ex int) {
	c := cli.pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		fmt.Println("set error", err.Error())
	}
}

func (cli *redisPool) CheckKey(k string) bool {
	c := cli.pool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return exist
	}
}

func (cli *redisPool) DelKey(k string) error {
	c := cli.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (cli *redisPool) SetJson(k string, data interface{}) error {
	c := cli.pool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func (cli *redisPool) GetJsonByte(key string) ([]byte, error) {
	c := cli.pool.Get()
	defer c.Close()
	jsonGet, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jsonGet, nil
}
