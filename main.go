package main

import (
	"ForumWeb/dao/mysql"
	"ForumWeb/dao/redis"
	"ForumWeb/log"
	"ForumWeb/pkg/kafka"
	"ForumWeb/pkg/snowflake"
	"ForumWeb/router"
	"ForumWeb/setting"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "f", "./conf/config.yaml", "配置信息文件路径")
	// 初始化配置信息
	if err := setting.Init(filepath); err != nil {
		fmt.Printf("init setting failed, err:%v\n", err)
		return
	}
	// 初始化日志
	if err := log.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init log failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	//初始化数据库
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	//启动redis的消息监听
	go redis.StartRedisPubSub()
	defer redis.Close()
	// 雪花生成ID算法初始化
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		zap.L().Error("init failed, err: %v\n", zap.Error(err))
		return
	}
	// 初始化kafka
	ctx := context.Background()
	go kafka.StartEmailConsumer(ctx)
	// 初始化路由
	r := router.Init(setting.Conf.Mode)
	// 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", setting.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个10秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个10秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 10秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过10秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
	// 多加一点东西，并发，kafka,websocket等等
	//分布式系统开发，负载均衡技术，系统容灾设计，高可用系统
}
