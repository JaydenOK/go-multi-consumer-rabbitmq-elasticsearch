package main

import (
	"app/libs/elasticsearchlib"
	"app/libs/mysqllib"
	"app/libs/redislib"
	"app/routers"
	"app/tasks"
	"app/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//加载配置
	loadConfig()

	//初始化存储服务
	mysqllib.InitMysqlClient()
	redislib.InitRedisClient()

	//初始化http路由
	r := gin.Default()
	routers.InitRouter(r)
	elasticsearchlib.InitESClient()
	//启动消费者监听协程
	tasks.Run()

	//启动服务
	httpPort := viper.GetString("app.httpPort")
	//1，启动Server
	//_ = r.Run(":" + httpPort)

	//2，使用http.Server内置的Shutdown()方法优雅地关机
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// 加载配置文件信息到viper
func loadConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(utils.StringToInterface(err.Error()))
	}
	fmt.Println("系统配置如下：")
	fmt.Println("app:", viper.Get("app"))
	fmt.Println("mysql:", viper.Get("mysql"))
	fmt.Println("redis:", viper.Get("redis"))
	fmt.Println("mongo:", viper.Get("mongo"))
	fmt.Println("rabbitmq:", viper.Get("rabbitmq"))
	fmt.Println("elasticsearch:", viper.Get("elasticsearch"))
}
