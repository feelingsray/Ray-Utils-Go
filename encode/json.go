package encode

import (
  "github.com/bytedance/sonic"
)

func LoadJson(in []byte, out any) error {
  if err := sonic.Unmarshal(in, out); err != nil {
    return err
  }
  return nil
}

func DumpJson(in any) ([]byte, error) {
  out, err := sonic.Marshal(in)
  if err != nil {
    return nil, err
  }
  return out, nil
}
