package encode

import "encoding/json"

func LoadJson(in []byte, out interface{}) error {
	err := json.Unmarshal(in, out)
	if err != nil {
		return err
	}
	return nil
}

/*
	对象转json
*/
func DumpJson(in interface{}) ([]byte, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}
