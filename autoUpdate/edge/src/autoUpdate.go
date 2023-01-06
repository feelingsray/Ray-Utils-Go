package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
	
	"gopkg.in/yaml.v2"
	
	"github.com/feelingsray/Ray-Utils-Go/encode"
	"github.com/feelingsray/Ray-Utils-Go/httpHelper"
	"github.com/feelingsray/Ray-Utils-Go/tools"
)

func EdgeDownload(conf *EdgeConf) error {
	// 下载releases.yaml文件
	DownloadRelease(conf)
	// 访问服务器获取发行清单
	release, err := releaseInfo(conf)
	if err != nil {
		return err
	}
	if release == nil {
		return errors.New("未获取到软件发行清单")
	}
	for _, r := range release.Previous {
		_, _ = DownloadFile(conf, r.Version, r.MD5)
	}
	return nil
}

func releaseInfo(conf *EdgeConf) (*ReleaseInfo, error) {
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

func DownloadRelease(conf *EdgeConf) {
	_, data, _ := httpHelper.HttpGet(fmt.Sprintf("http://%s/releases.yaml", conf.Remote))
	file, _ := os.OpenFile(path.Join(conf.Env.SourcePath, "releases.yaml"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer file.Close()
	_, _ = file.WriteString(data)
}

func DownloadFile(conf *EdgeConf, version string, fileMd5 string) (bool, error) {
	
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
	if ok, _ := tools.PathExists(path.Join(conf.Env.SourcePath, conf.Name)); !ok {
		_ = tools.CreateDir(path.Join(conf.Env.SourcePath, conf.Name))
	}
	
	client := http.DefaultClient
	client.Timeout = time.Second * 3600 //设置超时时间
	resp, err := client.Get(dUrl)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	
	// 获取文件名称
	filename := path.Base(uri.Path)
	// 下载临时文件
	tmpFilePath := path.Join(conf.Env.SourcePath, conf.Name, filename+".download")
	// 正式文件
	filePath := path.Join(conf.Env.SourcePath, conf.Name, filename)
	
	//创建临时文件
	file, err := os.Create(tmpFilePath)
	defer file.Close()
	
	if err != nil {
		return false, err
	}
	if resp.Body == nil {
		return false, errors.New("数据为空")
	}
	
	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return false, err
	}
	
	//io.copyBuffer() 的简化版本
	for {
		//读取bytes
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
