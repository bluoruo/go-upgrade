package lib

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// FileExists 文件是否存在
func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// FileCopy 复制文件
func FileCopy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		err = errors.New("不是一个常规文件！")
		return err
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	return nil
}

// FileReName 文件改名
func FileReName(name, dst string) bool {
	err := os.Rename(name, dst)
	if err != nil {
		fmt.Println("文件改名错误:", err)
		return false
	} else {
		return true
	}
}
