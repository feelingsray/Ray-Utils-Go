package commonIota

// Ledger 账本
type Ledger struct {
	DOID       string `json:"doid" bson:"doid"`
	DataHash   string `json:"data_hash" bson:"data_hash"`
	Data       any    `json:"data" bson:"data"`
	CsDataTime int64  `json:"cs_data_time" bson:"cs_data_time"`
}
