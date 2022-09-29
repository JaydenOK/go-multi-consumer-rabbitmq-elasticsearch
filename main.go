package main

import (
	"app/libs/elasticsearchlib"
	"app/libs/mysqllib"
	"app/libs/redislib"
	"app/routers"
	"app/tasks"
	"app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	_ = r.Run(":" + httpPort)
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
