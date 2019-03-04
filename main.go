package main

import (
	"flag"
	"logpro/defs"
	"logpro/monitor"
	"logpro/process"
	"logpro/read"
	"logpro/write"
	"time"
)

var (
	Path      string
	InfluxDsn string
)

func init() {
	flag.StringVar(&Path, "path", "D:/goproject/src/logpro/access.log", "log file path")
	flag.StringVar(&InfluxDsn, "influxDsn", "http://192.168.62.100:8086@admin@123456@log@s", "influx data source")
	flag.Parse()
}

func main() {
	lp := &proc.LogProcess{
		RC:    make(chan []byte),
		WC:    make(chan *defs.Message),
		Read:  &read.FromFile{Path: Path},
		Write: &write.ToInfluxDb{InfluxDbSn: InfluxDsn},
	}

	//从文件里读取日志内容
	go lp.Read.Read(lp.RC)

	//解析日志内容
	for i := 0; i < 2; i++ {
		go lp.Process() //解析模块比较慢，可以多加几个协程
	}

	//写入InfluxDB,写入如果是远程调用，可以多开几个，参数可以以传递的方式传入
	for i := 0; i < 4; i++ {
		go lp.Write.Write(lp.WC)
	}

	//系统状态监控
	m := &monitor.Monitor{
		StartTime: time.Now(),
		Data:      monitor.SystemInfo{},
	}
	m.Start(lp)
}
