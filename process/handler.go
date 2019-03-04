package proc

import (
	"log"
	"logpro/defs"
	"logpro/read"
	"logpro/write"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogProcess struct {
	RC    chan []byte
	WC    chan *defs.Message
	Read  read.Reader
	Write write.Writer
}

func (lp *LogProcess) Process() {

	/*
		日志内容：172.0.0.12 - - [04/mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
		正则表达式：([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	*/

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")

	for v := range lp.RC {
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail: ", string(v))
			continue
		}

		msg := &defs.Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail: ", err.Error(), ret[4])
			continue
		}
		msg.TimeLocal = t

		msg.BytesSent, _ = strconv.Atoi(ret[8])

		reqSlice := strings.Split(ret[6], " ")
		if len(reqSlice) != 3 {
			log.Println("strings split fail: ", ret[6])
			continue
		}
		msg.Method = reqSlice[0]

		u, err := url.Parse(reqSlice[1])
		if err != nil {
			log.Println("url parse fail: ", err)
		}
		msg.Path = u.Path
		msg.Scheme = ret[5]
		msg.Status = ret[7]
		msg.UpstreamTime, _ = strconv.ParseFloat(ret[12], 64)
		msg.RequestTime, _ = strconv.ParseFloat(ret[13], 64)

		lp.WC <- msg
	}
}
