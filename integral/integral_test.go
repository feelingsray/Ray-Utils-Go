package integral

import (
	"fmt"
	"testing"
)

func TestIntegral(t *testing.T) {
	item := new(Integral)
	if err := item.Init("192.168.111.142", 11531, "password01!", 20, DealReal, DealHis); err != nil {
		t.Fatal(err)
		return
	}
	static := &TagStatic{"tag01", "float", "-9999", "1",
		-9999, -9999, 40, 30,
		70, -9999, -9999, -9999}
	baseTimestamp := int64(1683877950)
	valueList := []float64{0, 15, 30, 45, 60, 75, 90, 115, 110, 101, 100, 85, 75, 65, 55, 10}

	for i, value := range valueList {
		newTag := new(TagReal)
		newTag.TagCode = "tag01"
		newTag.DataType = "float64"
		newTag.Value = value
		newTag.Timestamp = baseTimestamp + int64(i)
		err := item.Calculate("tag01", newTag, static)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func DealReal(alarmItem *AlarmItem) error {
	//fmt.Println(fmt.Sprintf("real:%+v", alarmItem))
	return nil
}

func DealHis(alarmItem *AlarmItem) error {
	if alarmItem.DuringType == "during" {
		return nil
	}
	fmt.Println(fmt.Sprintf("his:%+v", alarmItem))
	return nil
}
