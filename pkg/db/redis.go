package db

import (
	"encoding/json"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
)

var (
	ErrRedisNil = redis.ErrNil
)

type RedisConfig struct {
	Host   string
	Port   int
	Auth   string
	Select int
}

type Redis struct {
	pool *redis.Pool
}

func ConnRedis(conf RedisConfig) *Redis {
	r := &Redis{}
	r.pool = &redis.Pool{
		MaxIdle:     500,
		MaxActive:   100,
		IdleTimeout: time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", conf.Host+":"+strconv.Itoa(conf.Port),
				redis.DialPassword(conf.Auth),
				redis.DialDatabase(conf.Select),
				redis.DialConnectTimeout(2*time.Second),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second))
			if err != nil {
				panic("redis pool is nil")
			}
			return con, nil
		},
	}
	return r
}

type RedisCmd struct {
	reply interface{}
	err   error
}

func (r *Redis) Exists(key string) bool {
	conn := r.pool.Get()
	defer conn.Close()
	has, _ := redis.Bool(conn.Do("EXISTS", key))
	return has
}

func (r *Redis) Del(keys ...interface{}) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("DEL", keys...)
}

func (r *Redis) TTL(key string) int {
	conn := r.pool.Get()
	defer conn.Close()
	i, _ := redis.Int(conn.Do("TTL", key))
	return i
}

func (r *Redis) Get(key string) RedisCmd {
	return r.Do("GET", key)
}

func (r *Redis) SetEx(key string, ex time.Duration, val interface{}) error {
	return r.autoSet("SETEX", key, ex.Seconds(), val)
}

func (r *Redis) Incr(key string) int {
	conn := r.pool.Get()
	defer conn.Close()
	i, _ := redis.Int(conn.Do("INCR", key))
	return i
}

func (r *Redis) HGet(key string, field string) RedisCmd {
	return r.Do("HGET", key, field)
}

func (r *Redis) HSet(key string, field string, val interface{}) error {
	return r.autoSet("HSET", key, field, val)
}

func (r *Redis) autoSet(command string, args ...interface{}) error {
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

func (r *Redis) Do(commandName string, args ...interface{}) RedisCmd {
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
