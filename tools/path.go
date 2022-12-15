package tools

import (
  "io"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

func GetAppPath() string {
  file, _ := exec.LookPath(os.Args[0])
  path, _ := filepath.Abs(file)
  index := strings.LastIndex(path, string(os.PathSeparator))
  return path[:index]
}

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
  _, err = io.Copy(dst, src)
  if err != nil {
    return err
  }
  return nil
}

func GetAllFile(pathname string, s []string) ([]string, error) {
  rd, err := os.ReadDir(pathname)
  if err != nil {
    return s, err
  }
  for _, fi := range rd {
    if fi.IsDir() {
      fullDir := pathname + "/" + fi.Name()
      s, err = GetAllFile(fullDir, s)
      if err != nil {
        return s, err
      }
    } else {
      fullName := pathname + "/" + fi.Name()
      s = append(s, fullName)
    }
  }
  return s, nil
}
