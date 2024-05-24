package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/feelingsray/Ray-Utils-Go/encode"
	"github.com/feelingsray/Ray-Utils-Go/httpHelper"
	"github.com/feelingsray/Ray-Utils-Go/tarHelper"
	"github.com/feelingsray/Ray-Utils-Go/tools"
)

func CheckUpdate(conf *ClientConf, version string) error {
	// 访问服务器获取发行清单
	release, err := releaseInfo(conf)
	if err != nil {
		return err
	}
	if release == nil {
		return errors.New("未获取到软件发行清单")
	}

	softVer, err := softVersion(conf)
	if err != nil {
		return err
	}

	if softVer == nil || softVer["result"] == "err" {
		// 首次安装
		if version == "" {
			version = release.Latest
		}
		updateSoft(conf, version, release.Version[version])
	}
	if softVer != nil {
		if conf.Auto {
			// 自动升级为自动升级为最新版本
			if softVer["version"] == release.Latest && softVer["md5"] == release.MD5 {
			} else {
				updateSoft(conf, release.Latest, release.MD5)
			}
		} else {
			if version == "" {
				// 如果没有指定版本则升级为最新版本
				version = release.Latest
			}
			if softVer["version"] != version {
				// 升级或回退版本
				updateSoft(conf, version, release.Version[version])
			} else {
				if softVer["md5"] != release.Version[version] {
					// 当前版本MD5与清单不匹配,需要修正
					updateSoft(conf, version, release.Version[version])
				}
			}
		}
	}
	return nil
}

func updateSoft(conf *ClientConf, version string, md5Str string) {
	softPath := path.Join(conf.Env.FilePath, conf.Name, "release_"+version+".tar.gz")
	// 判断一下软件包是否存在,如果不存在则下载并解压
	if ok, _ := tools.PathExists(softPath); !ok {
		ok, err := DownloadFile(conf, version, md5Str)
		if err != nil {
			WriteSoftVersion(conf, version, md5Str, "err", "下载文件失败:"+err.Error())
			return
		}
		if !ok {
			WriteSoftVersion(conf, version, md5Str, "err", "下载文件失败")
			return
		}
		err = tarHelper.DeTarGzCompress(softPath, path.Join(conf.Env.FilePath, conf.Name)+"/")
		if err != nil {
			WriteSoftVersion(conf, version, md5Str, "err", "解压文件失败:"+err.Error())
			return
		}
	}

	if conf.Install {
		// 执行install.sh文件
		installPath := path.Join(conf.Env.FilePath, conf.Name, "release_"+version, "install.sh")
		exec.Command("/bin/bash", "-c", fmt.Sprintf("chmod +x %s", installPath)).Output()
		cmd := exec.Command("/bin/bash", "-c", installPath)
		bytes, err := cmd.Output()
		if err != nil {
			WriteSoftVersion(conf, version, md5Str, "err", "安装命令失败:"+err.Error())
			return
		}

		echoList := strings.Split(string(bytes), "\n")
		re := regexp.MustCompile("^.*success.*$")

		installFlag := false

		for i := 0; i < len(echoList); i++ {
			tmp := echoList[len(echoList)-i-1]
			if !re.MatchString(tmp) {
				installFlag = true
			}
		}

		if installFlag {
			WriteSoftVersion(conf, version, md5Str, "install", "安装成功")
			return
		} else {
			WriteSoftVersion(conf, version, md5Str, "err", "安装失败:"+string(bytes))
			return
		}

	} else {
		WriteSoftVersion(conf, version, md5Str, "down", fmt.Sprintf("下载文件成功:%s/release_%s.tar.gz", conf.Name, version))
		return
	}
}

func WriteSoftVersion(conf *ClientConf, version string, md5Str string, result string, msg string) {
	if ok, _ := tools.PathExists(conf.Env.VersionPath); !ok {
		_, _ = os.Create(conf.Env.VersionPath)
	}
	file, _ := os.OpenFile(conf.Env.VersionPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	defer file.Close()
	_, _ = file.WriteString(fmt.Sprintf("%s|%s|%d|%s|%s", version, md5Str, int(time.Now().Unix()), result, msg))
}

func releaseInfo(conf *ClientConf) (*ReleaseInfo, error) {
	_, data, _ := httpHelper.HttpGet(fmt.Sprintf("http://%s/releases.yaml", conf.Remote))
	releases := make([]*ReleaseInfo, 0)
	err := yaml.Unmarshal([]byte(data), &releases)
	if err != nil {
		return nil, err
	}
	for _, release := range releases {
		if release.Name == conf.Name {
			versionDict := make(map[string]string)
			for _, p := range release.Previous {
				if p.Version == release.Latest {
					release.MD5 = p.MD5
				}
				versionDict[p.Version] = p.MD5
			}
			release.Version = versionDict
			return release, nil
		}
	}
	return nil, nil
}

func softVersion(conf *ClientConf) (map[string]string, error) {
	if ok, _ := tools.PathExists(conf.Env.VersionPath); !ok {
		_, err := os.Create(conf.Env.VersionPath)
		return nil, err
	}
	fileObj, err := os.Open(conf.Env.VersionPath)
	if err != nil {
		return nil, err
	}
	defer fileObj.Close()
	content, err := io.ReadAll(fileObj)
	if err != nil {
		return nil, err
	}
	data := strings.Replace(string(content), " ", "", -1)
	if data == "" {
		// 首次写入
		return nil, nil
	}
	dataList := strings.Split(data, "|")
	result := make(map[string]string)
	result["version"] = dataList[0]
	result["md5"] = dataList[1]
	result["update"] = dataList[2]
	result["result"] = dataList[3]
	result["msg"] = dataList[4]
	return result, nil
}

func DownloadFile(conf *ClientConf, version string, fileMd5 string) (bool, error) {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	dUrl := fmt.Sprintf("http://%s/%s/release_%s.tar.gz/", conf.Remote, conf.Name, version)
	// 解析url
	uri, err := url.ParseRequestURI(dUrl)
	if err != nil {
		return false, err
	}
	// 下载文件目录是否存在
	if ok, _ := tools.PathExists(path.Join(conf.Env.FilePath, conf.Name)); !ok {
		_ = tools.CreateDir(path.Join(conf.Env.FilePath, conf.Name))
	}

	client := http.DefaultClient
	client.Timeout = time.Second * 3600 // 设置超时时间
	resp, err := client.Get(dUrl)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}

	// 获取文件名称
	filename := path.Base(uri.Path)
	// 下载临时文件
	tmpFilePath := path.Join(conf.Env.FilePath, conf.Name, filename+".download")
	// 正式文件
	filePath := path.Join(conf.Env.FilePath, conf.Name, filename)

	// 创建临时文件
	file, err := os.Create(tmpFilePath)
	defer file.Close()

	if err != nil {
		return false, err
	}
	if resp.Body == nil {
		return false, errors.New("数据为空")
	}

	// 读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return false, err
	}

	// io.copyBuffer() 的简化版本
	for {
		// 读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[:nr])
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if written >= fsize {
			// 写完数据了
			break
		}
	}

	if err != nil {
		return false, err
	}
	if written != fsize {
		return false, errors.New("文件下载不完整")
	}

	// 下面要操作文件了,所以先关闭了文件
	file.Close()

	if fileMd5 != encode.MD5File(tmpFilePath) {
		return false, errors.New("MD5校验失败:" + encode.MD5File(tmpFilePath))
	}

	// 如果文件存在先删除文件
	if ok, _ := tools.PathExists(filePath); ok {
		_ = os.Remove(filePath)
	}

	// 改名
	err = os.Rename(tmpFilePath, filePath)
	if err != nil {
		return false, nil
	}

	return true, nil
}
