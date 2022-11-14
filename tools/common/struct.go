package common

import (
	"encoding/json"
)

type KMessage struct {
	Head struct {
		Company        string `json:"company" bson:"company"`
		Biz            string `json:"biz" bson:"biz"`
		CSDataTime     int64  `json:"cs_data_time" bson:"cs_data_time"`
		ReceivedTime   int64  `json:"received_time" bson:"received_time"`
		GatewayCode    string `json:"gateway_code" bson:"gateway_code"`
		GatewayVersion string `json:"gateway_version" bson:"gateway_version"`
	} `json:"head" bson:"head"`
	Body string `json:"body" bson:"body"`
}

func Struct2Map(obj any) (map[string]any, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]any)
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}
