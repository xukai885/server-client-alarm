package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server-client-alarm/dao/mysql"
	"server-client-alarm/routes"
	"server-client-alarm/settings"
	"server-client-alarm/watchClient"
	"syscall"
	"time"
)

func main() {

	// 配置文件初始化
	if err := settings.Init(); err != nil {
		log.Println("配置文件初始化失败", err)
		return
	}

	// db初始化
	if err := mysql.Init(settings.Conf.Db); err != nil {
		log.Println("数据库初始化失败", err)
		return
	}

	defer mysql.Close()

	// 初始化读取数据库中id信息
	clientList, err := mysql.ClientInit()
	if err != nil {
		log.Println("mysql读取数据失败", err)
		return
	}

	watchClient.ClientListSum = clientList

	// 定时查询对象的信息
	go watchClient.Init()

	// 注册路由
	r := routes.SetUp()

	// 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	// 等待中断信号来优雅关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)                      // 创建一个接收信号的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM) //此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号是才会向下执行
	log.Println("shutdown server......")
	// 创建一个5秒的超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5秒内优雅关闭服务
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown:", err)
	}
	log.Println("server exiting")
}
