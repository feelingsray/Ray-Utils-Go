package src

import (
    "path"
    
    "github.com/feelingsray/Ray-Utils-Go/encode"
    "github.com/feelingsray/Ray-Utils-Go/tools"
)

type EdgeConf struct {
    Name   string `yaml:"name"`
    Remote string `yaml:"remote"`
    Env    struct {
        AppPath    string
        ConfPath   string
        SourcePath string
    }
}

func LoadEdgeConf() (*EdgeConf, error) {
    conf := EdgeConf{}
    appPath := tools.GetAppPath()
    appPath = "/Users/ray/jylink/Ray-Utils-Go/autoupdate/edge"
    confPath := path.Join(appPath, "autoupdate.yaml")
    err := encode.LoadYaml(confPath, &conf)
    if err != nil {
        return nil, err
    }
    conf.Env.SourcePath = path.Join(appPath, "releases")
    
    if ok, _ := tools.PathExists(conf.Env.SourcePath); !ok {
        err = tools.CreateDir(conf.Env.SourcePath)
        if err != nil {
            return nil, err
        }
    }
    conf.Env.AppPath = appPath
    conf.Env.ConfPath = confPath
    return &conf, nil
}

type ReleaseInfo struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Latest      string `yaml:"latest"`
    MD5         string
    Previous    []struct {
        Version string `yaml:"version"`
        MD5     string `yaml:"md5"`
    }
    Version map[string]string
}
