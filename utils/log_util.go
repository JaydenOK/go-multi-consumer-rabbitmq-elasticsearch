// Package name declaration
package utils

import (
	"fmt"
	"os"
)

// 日志记录 "./" 指当前工程目录下
const LogPath = "./logs/"

func WriteLog(str string) int {
	filePath := LogPath + "app-" + GetCurrentDate() + ".log"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("打开日志文件失败：", err.Error())
		return 0
	}
	//增加时间格式前缀及换行
	str = "[" + GetCurrentDateTime() + "]" + str + "\n"
	bytes := []byte(str)
	size, err := file.Write(bytes)
	if err != nil {
		fmt.Println("写入日志文件失败：", err.Error())
		return 0
	}
	return size
}
