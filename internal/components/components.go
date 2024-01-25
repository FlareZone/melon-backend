package components

import (
	"context"
	"fmt"
	"github.com/FlareZone/melon-backend/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

var (
	DBEngine *xorm.Engine
	Redis    *redis.Client
)

func initMysql() {
	// init DBEngine
	engine, err := xorm.NewEngine("mysql", config.MelonDB.DSN)
	if err != nil {
		panic(fmt.Errorf("connect to mysql fail, err is %v", err))
	}
	engine.ShowSQL(config.MelonDB.Logging)
	DBEngine = engine
}

func initRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisCfg.Addr,
		Password: config.RedisCfg.Password,
		DB:       config.RedisCfg.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Errorf("connect to redis fail, err is %v", err))
	}
	Redis = client
}

func InitComponents() {
	initMysql()
	initRedis()
}
