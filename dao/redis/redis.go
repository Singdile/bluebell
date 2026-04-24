package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// 声明全局的redis 客户端变量
var rdb *redis.Client

// 初始化redis客户端变量
func Init() (err error) {
	options := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: "",
		DB:       viper.GetInt("redis.db"),
		Protocol: viper.GetInt("redis.protocol"),
		PoolSize: viper.GetInt("redis.poolsize"),
	}

	rdb = redis.NewClient(&options)

	//测试连接
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("连接 redis 失败： %v\n", err)
		return err
	}

	fmt.Printf("连接成功\n")
	return nil
}

func Close() {
	_ = rdb.Close()
}
