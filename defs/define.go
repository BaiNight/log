package defs

import "time"

// 写入influxDB数据
type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// 用于统计处理成功、失败的数量
var TypeMonitorChan = make(chan int, 200)

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)
