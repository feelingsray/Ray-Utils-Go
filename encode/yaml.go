package encode

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
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
