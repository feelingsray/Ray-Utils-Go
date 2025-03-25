package integral

type AlarmItem struct {
	TagCode string

	// [dis,wval,aval,whi,wlow,ahi,alow,phi,plow,rhi,rlow]
	AlarmType string // 报警类型

	DataType   string // 数据类型
	DuringType string // 报警持续类型[start, during, end]
	Threshold  any    // 报警线：开关量报警值/模拟量阈值/其他报警指标

	// 所有iot报警都有的属性
	StartVal any // 报警开始值

	StartTime   int64 // 报警开始时间
	LastTagTime int64 // 上一个报警点时间

	// 下面六个属性，只有积分报警用到
	Integral float64 // 当前积分总量

	LastTagValue float64 // 最后一个点的取值
	MaxValue     float64 // 最大值
	MinValue     float64 // 最小值
	MaxTime      int64   // 报警最大值时间
	MinTime      int64   // 报警最小值时间
}

type TagReal struct {
	DataType  string // 数据类型
	TagCode   string // 测点编码
	Value     any    // 传感器数值
	Timestamp int64  // 传感器测量时间
}

type TagStatic struct {
	TagCode  string
	DataType string
	WVal     string
	AVal     string
	WHI      float64
	WLOW     float64
	AHI      float64
	ALOW     float64
	PHI      float64
	PLOW     float64
	RHI      float64
	RLOW     float64
}
