package modbus

type ReadRes struct {
	Write []byte
	Read  []byte
	Val   float64
}
