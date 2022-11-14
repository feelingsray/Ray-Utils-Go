package common

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
