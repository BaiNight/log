package write

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"logpro/defs"
	"strings"
)

type Writer interface {
	Write(<-chan *defs.Message)
}

type ToInfluxDb struct {
	InfluxDbSn string
}

func (w *ToInfluxDb) Write(wc <-chan *defs.Message) {
	infSli := strings.Split(w.InfluxDbSn, "@")

	// 新建一个客户端连接
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  infSli[3],
		Precision: infSli[4], //精度
	})
	if err != nil {
		log.Fatal(err)
	}

	for v := range wc {
		tags := map[string]string{
			"Path":   v.Path,
			"Method": v.Method,
			"Scheme": v.Scheme,
			"Status": v.Status,
		}
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// 批量写入
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// 关闭客户端资源
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("write successfully")
	}
}
