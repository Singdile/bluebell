package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// @title BlueBell 社区论坛 API
// @version 1.0
// @description bluebell社区论坛后端API

// @contact.name singdile
// @contact.url https://github.com/singdile
// @contact.email singdile0709@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

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

	// 5.同步 Redis 索引（从 MySQL 全量同步）
	ctx := context.Background()
	if err := redis.SyncPostsFromMySQL(ctx); err != nil {
		zap.L().Error("failed to sync redis index", zap.Error(err))
		// 不阻止服务启动，只是记录错误
	}
	zap.L().Info("sync redis index success")

	// 6.初始化 score 更新队列
	mysql.InitScoreQueue()
	defer mysql.StopScoreQueue()
	zap.L().Info("init score queue success")

	// 初始化 vote 更新队列
	mysql.InitVoteUpdateQueue()
	defer mysql.StopVoteUpdateQueue()
	zap.L().Info("init vote queue success")


	// 7.初始化雪花算法计算机节点
	fmt.Printf("time: %v %v\n", settings.Conf.Snowflakeconfig.StartTime, settings.Conf.Snowflakeconfig.MachineID)
	if err := snowflake.Init(settings.Conf.Snowflakeconfig.StartTime, settings.Conf.Snowflakeconfig.MachineID); err != nil {
		zap.L().Error("failed to init snowflake node", zap.Error(err))
		return
	}
	zap.L().Info("init snowflake node success")

	// 8.注册路由
	r := router.SetupRouter()

	// 9.启动服务
	r.Run()

}
