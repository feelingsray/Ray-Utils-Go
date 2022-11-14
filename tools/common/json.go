package common

import (
	"github.com/json-iterator/go"
)

func LoadJson(in []byte, out any) error {
	err := jsoniter.Unmarshal(in, out)
	if err != nil {
		return err
	}
	return nil
}

func DumpJson(in any) ([]byte, error) {
	out, err := jsoniter.Marshal(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}
