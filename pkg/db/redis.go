package db

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var (
	Redis myRedis

	ErrRedisNil = redis.ErrNil
)

// redis
type redisConfig struct {
	Host   string
	Port   int
	Auth   string
	Select int
}

type myRedis struct {
	pool *redis.Pool
}

func ConnectRedis() {
	var conf redisConfig
	viper.UnmarshalKey("redis", &conf)

	conn, err := redis.Dial("tcp", conf.Host+":"+strconv.Itoa(conf.Port),
		redis.DialPassword(conf.Auth),
		redis.DialDatabase(conf.Select),
		redis.DialConnectTimeout(time.Second*2),
	)
	if err != nil {
		panic(err)
	}

	Redis.pool = &redis.Pool{
		IdleTimeout: time.Second * 2,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	fmt.Println("Redis连接成功")
}

type RedisCmd struct {
	reply interface{}
	err   error
}

func (r *myRedis) Exists(key string) bool {
	conn := r.pool.Get()
	defer conn.Close()
	has, _ := redis.Bool(conn.Do("EXISTS", key))
	return has
}

func (r *myRedis) Del(keys ...interface{}) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("DEL", keys...)
}

func (r *myRedis) TTL(key string) int {
	conn := r.pool.Get()
	defer conn.Close()
	i, _ := redis.Int(conn.Do("TTL", key))
	return i
}

func (r *myRedis) Get(key string) RedisCmd {
	return r.Do("GET", key)
}

func (r *myRedis) SetEx(key string, ex time.Duration, val interface{}) error {
	return r.autoSet("SETEX", key, ex.Seconds(), val)
}

func (r *myRedis) Incr(key string) int {
	conn := r.pool.Get()
	defer conn.Close()
	i, _ := redis.Int(conn.Do("INCR", key))
	return i
}

func (r *myRedis) HGet(key string, field string) RedisCmd {
	return r.Do("HGET", key, field)
}

func (r *myRedis) HSet(key string, field string, val interface{}) error {
	return r.autoSet("HSET", key, field, val)
}

func (r *myRedis) autoSet(command string, args ...interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	val := args[len(args)-1]
	switch val.(type) {
	case string, int, uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, bool:
		_, err := conn.Do(command, args...)
		return err
	default:
		buff, err := json.Marshal(val)
		if err == nil {
			args[len(args)-1] = string(buff)
			_, err = conn.Do(command, args...)
		}
		return err
	}
}

func (r *myRedis) Do(commandName string, args ...interface{}) RedisCmd {
	conn := r.pool.Get()
	defer conn.Close()
	cmd := RedisCmd{}
	cmd.reply, cmd.err = conn.Do(commandName, args...)
	return cmd
}

func (rc RedisCmd) Int() (int, error) {
	return redis.Int(rc.reply, rc.err)
}

func (rc RedisCmd) String() (string, error) {
	return redis.String(rc.reply, rc.err)
}

func (rc RedisCmd) Bool() (bool, error) {
	return redis.Bool(rc.reply, rc.err)
}

func (rc RedisCmd) Scan(to interface{}) error {
	if rc.err != nil {
		return rc.err
	}
	b, _ := redis.Bytes(rc.reply, nil)
	json.Unmarshal(b, to)
	return nil
}
