package src

import (
    "path"
    
    _ "gopkg.in/yaml.v2"
    
    "github.com/feelingsray/Ray-Utils-Go/encode"
    "github.com/feelingsray/Ray-Utils-Go/tools"
)

type ClientConf struct {
    Name    string `yaml:"name"`
    Remote  string `yaml:"remote"`
    Auto    bool   `yaml:"auto"`
    Install bool   `yaml:"install"`
    Source  string `yaml:"source"`
    Env     struct {
        AppPath     string
        ConfPath    string
        FilePath    string
        VersionPath string
    }
}

func LoadClientConf() (*ClientConf, error) {
    conf := ClientConf{}
    appPath := tools.GetAppPath()
    appPath = "/Users/ray/jylink/Ray-Utils-Go/autoupdate/client"
    confPath := path.Join(appPath, "autoupdate.yaml")
    err := encode.LoadYaml(confPath, &conf)
    if err != nil {
        return nil, err
    }
    if conf.Source == "" {
        conf.Source = path.Join(appPath, "release")
    }
    if ok, _ := tools.PathExists(conf.Source); !ok {
        err = tools.CreateDir(conf.Source)
        if err != nil {
            return nil, err
        }
    }
    conf.Env.AppPath = appPath
    conf.Env.ConfPath = confPath
    conf.Env.FilePath = conf.Source
    conf.Env.VersionPath = path.Join(appPath, "version")
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
