package main

import (
	"app/lib/mysqllib"
	"app/lib/redislib"
	"app/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	loadConfig()
	mysqllib.InitMysqlClient()
	redislib.InitRedisClient()

	r := gin.Default()
	routers.InitRouter(r)

	httpPort := viper.GetString("app.httpPort")
	r.Run(":" + string(httpPort))
}

// 加载配置文件信息到viper
func loadConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println("app:", viper.Get("app"))
	fmt.Println("mysql:", viper.Get("mysql"))
	fmt.Println("redis:", viper.Get("redis"))
	fmt.Println("mongo:", viper.Get("mongo"))
	fmt.Println("rabbitmq:", viper.Get("rabbitmq"))
	fmt.Println("elasticsearch:", viper.Get("elasticsearch"))
}
