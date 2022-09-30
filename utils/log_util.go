// Package name declaration
package utils

import (
	"fmt"
	"os"
)

const LogPath = "../logs/"

func WriteLog(str string) int {
	filePath := LogPath + GetCurrentDate() + "-" + ".txt"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("打开日志文件失败：", err.Error())
		return 0
	}
	bytes := []byte(str)
	size, err := file.Write(bytes)
	if err != nil {
		fmt.Println("写入日志文件失败：", err.Error())
		return 0
	}
	return size
}
