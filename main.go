package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"fmt"
	"go.uber.org/zap"
)

// go web 开发常用的脚手架模板
func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		panic(err)
	}
	fmt.Printf("config init success\n")

	// 2.初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed,err : %v\n", err)
	}
	defer zap.L().Sync()
	zap.L().Info("logger init success")

	// 3.初始化数据库连接(mysql)
	if err := mysql.Init(); err != nil {
		zap.L().Error("failed to init mysql", zap.Any("error", err))
		return
	}
	defer mysql.Close()
	zap.L().Info("init mysql success")

	// 4.初始化Redis
	if err := redis.Init(); err != nil {
		zap.L().Error("failed to init redis", zap.Any(("error"), err))
		return
	}
	defer redis.Close()
	zap.L().Info("init redis success")

	// 5.初始化雪花算法计算机节点
	fmt.Printf("time: %v %v\n", settings.Conf.Snowflakeconfig.StartTime, settings.Conf.Snowflakeconfig.MachineID)
	if err := snowflake.Init(settings.Conf.Snowflakeconfig.StartTime, settings.Conf.Snowflakeconfig.MachineID); err != nil {
		zap.L().Error("failed to init snowflake node", zap.Error(err))
		return
	}
	zap.L().Info("init snowflake node success")

	// 5.注册路由
	r := router.SetupRouter()

	// 6.启动服务
	r.Run()

}
