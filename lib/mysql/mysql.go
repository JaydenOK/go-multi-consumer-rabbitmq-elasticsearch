package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlClient *gorm.DB
)

func InitMysqlClient() {
	var dsn string
	host := viper.GetString("mysql.host")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	port := viper.GetString("mysql.port")
	database := viper.GetString("mysql.database")
	//连接数据库的时候加入参数parseTime=true 和loc=Local ，解决时间格式化问题
	dsn = username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql连接异常：", err.Error())
		panic(err)
	}
	mysqlClient = db
}

func GetMysqlClient() *gorm.DB {
	return mysqlClient
}
