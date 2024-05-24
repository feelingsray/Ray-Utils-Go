package encode

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYaml(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
