package logger

import (
	"bluebell/settings"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 自定义日志配置
func Init() (err error) {
	//将编码的时间格式指定为人类用户友好的格式
	encodeconfig := zap.NewProductionEncoderConfig()
	encodeconfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encodeconfig)                                                     //如何写入文件
	file, _ := os.OpenFile(settings.Conf.Logconfig.Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644) //写入什么日志文件
	writesyncer := zapcore.AddSync(file)

	var l = new(zapcore.Level)
	level := settings.Conf.Logconfig.Level
	err = l.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}

	core := zapcore.NewCore(encoder, writesyncer, zapcore.DebugLevel)

	lg := zap.New(core, zap.AddCaller()) //选项显示调用的函数

	//zap.AddCaller 指定显示调用方函数的信息
	zap.ReplaceGlobals(lg) // 替换掉zap实例中维护的全局的  zap.Logger, 在其他地方可以通过zap.L()访问到
	fmt.Printf("logger init success\n")
	return nil
}

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		//执行后续处理逻辑
		c.Next()

		//后续处理逻辑结束后,收集信息
		cost := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		}

		//记录信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error_msg", c.Errors.ByType(gin.ErrorTypeAny).String()))
			zap.L().Error(path, fields...)
		} else {
			zap.L().Info(path, fields...)
		}

	}
}

// ZapRecovery 是一个 Gin 中间件，用于替代官方默认的 Recovery。
// 它能捕获 panic，利用 zap 记录结构化日志，并根据连接状态优雅返回响应。
func ZapRecovery(stack bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			//尝试从panic中恢复，如果rex为nil,说明没有panic，正常退出
			if rec := recover(); rec != nil {
				var isBrokenPipe bool //检查连接是否已经中断

				err, ok := rec.(error) //类型断言，查看panic抛出的是不是error类型

				//查看error对应是否是这几种error
				//1.syscall.EPIPE 管道破裂,用户关闭了连接
				//2.syscall.ECONNECTSET 连接重置
				//3.http.ErrAbortHandler 处理函数主动停止处理并丢弃这个请求
				if ok {
					isBrokenPipe = errors.Is(err, syscall.EPIPE) ||
						errors.Is(err, syscall.ECONNRESET) ||
						errors.Is(err, http.ErrAbortHandler)
				}

				//获取请求快照，记录用户发送了什么包导致了panic
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)

				if zap.L() != nil {
					//构造日志字段
					fields := []zap.Field{
						zap.Any("error", rec),
						zap.String("request", string(httpRequest)),
					}

					//根据参数决定是否要添加堆栈信息
					if stack && !isBrokenPipe {
						fields = append(fields, zap.String("stack", string(debug.Stack())))
					}

					//日志记录错误
					if isBrokenPipe {
						zap.L().Error(ctx.Request.URL.Path, fields...)
					} else {
						zap.L().Error("recovery form panic", fields...)
					}
				}

				//响应控制
				//如果是连接中断,记录错误并中断后续的执行
				if isBrokenPipe {
					_ = ctx.Error(err.(error))
					ctx.Abort()
				} else {
					//业务崩溃,但是还有连接,返回500错误表示
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"code":    500,
						"message": "服务器内部错误",
					})
				}
			}
		}()

		//将请求交给下一个中间件处理
		ctx.Next()
	}
}
