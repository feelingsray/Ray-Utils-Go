package encode

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
	载入yaml格式文件并解析
*/
func LoadYaml(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
