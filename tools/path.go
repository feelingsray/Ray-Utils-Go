package tools

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 获取系统路径
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

// 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dir string) error {
	isExist, err := PathExists(dir)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	} else {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func CopyFile(dstName string, srcName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer dst.Close()
	_,err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
