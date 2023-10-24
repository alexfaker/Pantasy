package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

type RedisBusiness string //业务类型

const (
	RedisBusinessFirst  = RedisBusiness("first")
	RedisBusinessSecond = RedisBusiness("second")
)

// RedisGet 获取一个Redis实例
func RedisGet(name RedisBusiness, db int) *redis.Client {
	return localRedisManager.get(name, db)
}

var localRedisManager = &RedisDatabaseManager{
	Databases: map[RedisBusiness]*redisConfig{},
}

type redisConfig struct {
	Address  string   `yaml:"address"`  //地址
	Password string   `yaml:"password"` //密码
	clients  sync.Map `yaml:"-"`        //客户端列表
}

type RedisDatabaseManager struct {
	Databases map[RedisBusiness]*redisConfig
}

func (t *RedisDatabaseManager) add(name RedisBusiness, address, password string) {
	t.Databases[name] = &redisConfig{
		Address:  address,
		Password: password,
		clients:  sync.Map{},
	}
}

func (t *RedisDatabaseManager) get(name RedisBusiness, db int) *redis.Client {
	rdb, ok := t.Databases[name]
	if ok == false {
		return nil
	}
	if cli, ok := rdb.clients.Load(db); ok {
		return cli.(*redis.Client)
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     rdb.Address,
		Password: rdb.Password,
		DB:       db,
	})
	rdb.clients.Store(db, cli)

	return cli
}

func (t *RedisDatabaseManager) Ping() {
	for k, kVal := range t.Databases {
		db := t.get(k, 0)
		pong, err := db.Ping(context.Background()).Result()
		if err != nil {
			xlog.Fatalf("Redis ping %s, pong: %+v, err: %+v", kVal.Address, pong, err)
			return
		}
	}
}
