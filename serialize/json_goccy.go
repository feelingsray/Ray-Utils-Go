//go:build !std

package serialize

import (
	"github.com/goccy/go-json"
)

func LoadJson(in []byte, out any) error {
	if err := json.Unmarshal(in, out); err != nil {
		return err
	}
	return nil
}

func DumpJson(in any) ([]byte, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}
