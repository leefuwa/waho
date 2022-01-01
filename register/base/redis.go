package base

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
	"waho/register"
)

type RedisRegister struct {
	register.BaseRegister
}

var rdb *redis.Client
var ctx = context.Background()

func (redisRegister *RedisRegister) Init() {
	_, _ = NewRdb(register.GetConf.Section("redis").GetIntDefault("name", 0))
}

func Rdb() *redis.Client {
	if rdb == nil {
		NewRdb(register.GetConf.Section("redis").GetIntDefault("name", 0))
	}
	return rdb
}

func NewRdb(db int) (*redis.Client, error) {
	host := register.GetConf.Section("redis").GetDefault("host", "127.0.0.1")
	port := register.GetConf.Section("redis").GetDefault("port", "6379")
	pwd := register.GetConf.Section("redis").GetDefault("pwd", "")
	rdb = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panic("redis连接失败，", err)
		return nil, errors.New("redis连接失败，" + err.Error())
	}

	return rdb, nil
}
func RdbSet(key string, val interface{}, expiration time.Duration) *redis.StatusCmd {
	return Rdb().Set(ctx, key, val, expiration)
}
func RdbGet(key string) (string, error)  {
	return Rdb().Get(ctx, key).Result()
}
func RdbHSet(key string, values ...interface{}) *redis.IntCmd  {
	return Rdb().HSet(ctx, key, values...)
}
func RdbHGet(key string, field string) (string, error)  {
	return Rdb().HGet(ctx, key, field).Result()
}
func RdbHGetAll(key string) (map[string]string, error)  {
	return Rdb().HGetAll(ctx, key).Result()
}
func RdbHDel(key string, field string) (int64, error) {
	return Rdb().HDel(ctx, key, field).Result()
}
func RdbExpire(key string, expiration time.Duration)  {
	Rdb().Expire(ctx, key, expiration)
}

func RdbExpireAt(key string, tm time.Time)  {
	Rdb().ExpireAt(ctx, key, tm)
}


