package file

import (
	"log"
	"os"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}

// CreateFileWithParentDir 如果文件目录不存在，直接创建会失败，需要先创建目录
func CreateFileWithParentDir(path string) {
	index := strings.LastIndex(path, "/")
	if index > 0 {
		dirPath := path[0:index]
		if len(dirPath) > 0 && !IsExist(dirPath) {
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if !IsExist(path) {
		_, err := os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
