package integral

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"github.com/feelingsray/ray-utils-go/v2/serialize"
)

type Integral struct {
	RedisConf struct {
		Host     string
		Port     int
		Password string
		DBNum    int
	}
	redisStore []*redis.Client
	realFunc   dealRealAlarmItem
	hisFunc    dealHisAlarmItem
}

type (
	dealRealAlarmItem func(alarmItem *AlarmItem) error
	dealHisAlarmItem  func(alarmItem *AlarmItem) error
)

func (integral *Integral) Init(host string, port int, password string, dbNum int, realFunc dealRealAlarmItem, hisFunc dealHisAlarmItem) error {
	integral.RedisConf.Host = host
	integral.RedisConf.Port = port
	integral.RedisConf.Password = password
	integral.RedisConf.DBNum = dbNum
	integral.realFunc = realFunc
	integral.hisFunc = hisFunc

	for i := 0; i < integral.RedisConf.DBNum; i++ {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", integral.RedisConf.Host, integral.RedisConf.Port),
			Password: integral.RedisConf.Password,
			DB:       i,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return fmt.Errorf("redis connect err: %s", err.Error())
		} else {
			integral.redisStore = append(integral.redisStore, rdb)
		}
	}
	return nil
}

func (integral *Integral) Clear() error {
	if len(integral.redisStore) > 0 {
		cmdStatus := integral.redisStore[0].FlushAll(context.Background())
		if cmdStatus.Err() != nil {
			return fmt.Errorf("redis flushdb err: %s", cmdStatus.Err())
		}
	}
	return nil
}

func (integral *Integral) Calculate(tagCode string, newTag *TagReal, dbTag *TagStatic) error {
	if ((dbTag.WHI == dbTag.WLOW) && dbTag.WHI != -9999) || ((dbTag.AHI == dbTag.ALOW) && dbTag.AHI != -9999) ||
		((dbTag.PHI == dbTag.PLOW) && dbTag.PHI != -9999) || ((dbTag.RHI == dbTag.RLOW) && dbTag.RHI != -9999) {
		return nil
	}
	if newTag.Value == nil {
		return nil
	}

	hasAlarm := false
	var alarmItem *AlarmItem
	ntv := cast.ToFloat64(newTag.Value)
	if ntv == -9999 || ntv == -999 || (fmt.Sprintf("%v", newTag.Value) == "-9999" || fmt.Sprintf("%v", newTag.Value) == "-999") {
		alarmItem, _ = integral.dealAlarmTag(tagCode, "dis", newTag, dbTag)
	} else {
		alarmItem, _ = integral.dealNormalTag(tagCode, "dis", newTag)
	}
	if alarmItem != nil {
		hasAlarm = true
		_ = integral.hisFunc(alarmItem)
	}
	// 处理实时报警，并且不继续下面的计算
	if hasAlarm {
		_ = integral.realFunc(alarmItem)
		return nil
	}

	// 多态量和开关量预警
	alarmItem = nil
	if dbTag.WVal != "" && (dbTag.DataType == "bool" || dbTag.DataType == "string") {
		if fmt.Sprintf("%v", newTag.Value) == dbTag.WVal {
			alarmItem, _ = integral.dealAlarmTag(tagCode, "wval", newTag, dbTag)
		} else {
			alarmItem, _ = integral.dealNormalTag(tagCode, "wval", newTag)
		}
	}
	if alarmItem != nil {
		hasAlarm = true
		_ = integral.hisFunc(alarmItem)
	}
	// 处理实时报警，并且不继续下面的计算
	if hasAlarm {
		_ = integral.realFunc(alarmItem)
		return nil
	}

	// 多态量和开关量报警
	alarmItem = nil
	if dbTag.AVal != "" && (dbTag.DataType == "bool" || dbTag.DataType == "string") {
		if fmt.Sprintf("%v", newTag.Value) == dbTag.WVal {
			alarmItem, _ = integral.dealAlarmTag(tagCode, "aval", newTag, dbTag)
		} else {
			alarmItem, _ = integral.dealNormalTag(tagCode, "aval", newTag)
		}
	}
	if alarmItem != nil {
		hasAlarm = true
		_ = integral.hisFunc(alarmItem)
	}
	// 处理实时报警，并且不继续下面的计算
	if hasAlarm {
		_ = integral.realFunc(alarmItem)
		return nil
	}

	// 模拟量报警
	var alarmItemTmp *AlarmItem
	if dbTag.DataType == "float" || dbTag.DataType == "int" {
		// 先把int64统一变为float64
		if reflect.TypeOf(newTag.Value).String() == "int" {
			newTag.Value = float64(newTag.Value.(int))
		}
		// 是int64或者float64
		if reflect.TypeOf(newTag.Value).String() == "float64" {
			if dbTag.WHI == dbTag.WLOW {
				alarmItemTmp = nil
				if newTag.Value.(float64) > dbTag.WHI && dbTag.WHI != -9999 {
					// 超预警上限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "whi", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "whi", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
				alarmItemTmp = nil
				if newTag.Value.(float64) < dbTag.WLOW && dbTag.WLOW != -9999 {
					// 超预警下限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "wlow", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "wlow", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
			}

			if dbTag.AHI != dbTag.ALOW {
				alarmItemTmp = nil
				if newTag.Value.(float64) > dbTag.AHI && dbTag.AHI != -9999 {
					// 超报警上限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "ahi", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "ahi", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
				alarmItemTmp = nil
				if newTag.Value.(float64) < dbTag.ALOW && dbTag.ALOW != -9999 {
					// 超报警下限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "alow", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "alow", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
			}

			if dbTag.PHI != dbTag.PLOW {
				alarmItemTmp = nil
				if newTag.Value.(float64) > dbTag.PHI && dbTag.PHI != -9999 {
					// 超断电上限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "phi", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "phi", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
				alarmItemTmp = nil
				if newTag.Value.(float64) < dbTag.PLOW && dbTag.PLOW != -9999 {
					// 超断电下限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "plow", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "plow", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
			}

			if dbTag.RHI != dbTag.RLOW {
				alarmItemTmp = nil
				if newTag.Value.(float64) > dbTag.RHI && dbTag.RHI != -9999 {
					// 超量程上限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "rhi", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "rhi", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
				alarmItemTmp = nil
				if newTag.Value.(float64) < dbTag.RLOW && dbTag.RLOW != -9999 {
					// 超量程下限
					alarmItemTmp, _ = integral.dealAlarmTag(tagCode, "rlow", newTag, dbTag)
				} else {
					alarmItemTmp, _ = integral.dealNormalTag(tagCode, "rlow", newTag)
				}
				if alarmItemTmp != nil {
					hasAlarm = true
					alarmItem = alarmItemTmp
					_ = integral.hisFunc(alarmItem)
				}
			}
		}
	}
	// 模拟量的实时报警，总共处理一次
	if hasAlarm {
		_ = integral.realFunc(alarmItem)
	}
	return nil
}

func (integral *Integral) dealAlarmTag(tagCode string, alarmType string,
	newTag *TagReal, dbTag *TagStatic,
) (*AlarmItem, error) {
	alarmItem, err := integral.getAlarmItem(tagCode, alarmType)
	if err != nil {
		return nil, nil
	}
	if alarmItem == nil {
		newAlarmItem := new(AlarmItem)
		newAlarmItem.TagCode = tagCode
		newAlarmItem.AlarmType = alarmType
		newAlarmItem.DuringType = "start"
		switch alarmType {
		case "dis":
			newAlarmItem.Threshold = "-9999"
			newAlarmItem.DataType = "dis"
		case "wval":
			newAlarmItem.Threshold = dbTag.AVal
			newAlarmItem.DataType = "bool"
		case "aval":
			newAlarmItem.Threshold = dbTag.AVal
			newAlarmItem.DataType = "bool"
		case "whi":
			newAlarmItem.Threshold = dbTag.WHI
			newAlarmItem.DataType = "float"
		case "wlow":
			newAlarmItem.Threshold = dbTag.WLOW
			newAlarmItem.DataType = "float"
		case "ahi":
			newAlarmItem.Threshold = dbTag.AHI
			newAlarmItem.DataType = "float"
		case "alow":
			newAlarmItem.Threshold = dbTag.ALOW
			newAlarmItem.DataType = "float"
		case "phi":
			newAlarmItem.Threshold = dbTag.PHI
			newAlarmItem.DataType = "float"
		case "plow":
			newAlarmItem.Threshold = dbTag.PLOW
			newAlarmItem.DataType = "float"
		case "rhi":
			newAlarmItem.Threshold = dbTag.RHI
			newAlarmItem.DataType = "float"
		case "rlow":
			newAlarmItem.Threshold = dbTag.RLOW
			newAlarmItem.DataType = "float"
		default:
			return nil, nil
		}
		newAlarmItem.StartVal = newTag.Value
		newAlarmItem.StartTime = newTag.Timestamp
		newAlarmItem.LastTagTime = newTag.Timestamp
		newAlarmItem.Integral = 0
		if newAlarmItem.DataType == "float" {
			newValue, err := strconv.ParseFloat(fmt.Sprintf("%v", newTag.Value), 64)
			if err != nil {
				return nil, nil
			}
			newAlarmItem.LastTagValue = newValue
			newAlarmItem.MaxValue = newValue
			newAlarmItem.MinValue = newValue
			newAlarmItem.MaxTime = newTag.Timestamp
			newAlarmItem.MinTime = newTag.Timestamp
		}
		if err = integral.setAlarmItem(tagCode, alarmType, newAlarmItem); err != nil {
			return nil, err
		}
		return newAlarmItem, nil
	} else {
		if alarmItem.LastTagTime > newTag.Timestamp {
			return nil, nil
		}
		alarmItem.DuringType = "during"
		if alarmItem.DataType == "float" {
			newValue, err := cast.ToFloat64E(newTag.Value)
			if err != nil {
				return nil, nil
			}
			alarmItem.Integral += (newValue + alarmItem.LastTagValue) * float64(newTag.Timestamp-alarmItem.LastTagTime) / 2
			alarmItem.LastTagValue = newValue
			if newValue > alarmItem.MaxValue {
				alarmItem.MaxValue = newValue
				alarmItem.MaxTime = newTag.Timestamp
			}
			if newValue < alarmItem.MinValue {
				alarmItem.MinValue = newValue
				alarmItem.MinTime = newTag.Timestamp
			}
		}
		alarmItem.LastTagTime = newTag.Timestamp
		_ = integral.setAlarmItem(tagCode, alarmType, alarmItem)
		return alarmItem, nil
	}
}

func (integral *Integral) dealNormalTag(tagCode string, alarmType string,
	newTag *TagReal,
) (*AlarmItem, error) {
	alarmItem, err := integral.getAlarmItem(tagCode, alarmType)
	if err != nil {
		return nil, err
	}
	if alarmItem != nil {
		if alarmItem.LastTagTime > newTag.Timestamp {
			return alarmItem, nil
		}
		_ = integral.deleteAlarmItem(tagCode, alarmType)
		alarmItem.DuringType = "end"
		newValue, err := cast.ToFloat64E(newTag.Value)
		if err != nil {
			return alarmItem, nil
		}
		alarmItem.Integral += (newValue + alarmItem.LastTagValue) * float64(newTag.Timestamp-alarmItem.LastTagTime) / 2
		alarmItem.LastTagTime = newTag.Timestamp
		alarmItem.LastTagValue = cast.ToFloat64(newTag.Value)
		return alarmItem, nil
	}
	return nil, nil
}

func (integral *Integral) setAlarmItem(tagCode string, alarmType string, alarmItem *AlarmItem) error {
	alarmItemByte, err := serialize.DumpJson(*alarmItem)
	if err != nil {
		return fmt.Errorf("json dump err: %s", err.Error())
	}
	route := TakeMold(tagCode, integral.RedisConf.DBNum)
	err = integral.redisStore[route].Set(
		context.Background(),
		fmt.Sprintf("%s_%s", tagCode, alarmType),
		string(alarmItemByte),
		-1,
	).Err()
	if err != nil {
		return fmt.Errorf("redis set err: %s", err.Error())
	}
	return nil
}

func (integral *Integral) getAlarmItem(tagCode string, alarmType string) (*AlarmItem, error) {
	route := TakeMold(tagCode, integral.RedisConf.DBNum)
	alarmString, err := integral.redisStore[route].Get(
		context.Background(),
		fmt.Sprintf("%s_%s", tagCode, alarmType),
	).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis get err: %s", err.Error())
	}
	if len(alarmString) == 0 {
		return nil, nil
	} else {
		alarmItem := AlarmItem{}
		err = serialize.LoadJson([]byte(alarmString), &alarmItem)
		if err != nil {
			return nil, fmt.Errorf("json load err: %s", err.Error())
		}
		return &alarmItem, nil
	}
}

func (integral *Integral) deleteAlarmItem(tagCode string, alarmType string) error {
	route := TakeMold(tagCode, integral.RedisConf.DBNum)
	err := integral.redisStore[route].Del(
		context.Background(),
		fmt.Sprintf("%s_%s", tagCode, alarmType),
	).Err()
	if err != nil {
		return fmt.Errorf("redis删除报警对象失败:%s", err.Error())
	}
	return nil
}
