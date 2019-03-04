package monitor

import (
	"encoding/json"
	"io"
	"logpro/process"
	"net/http"
	"time"
)

type Monitor struct {
	StartTime time.Time
	Data      SystemInfo
	TpsSlic   []int
}

//系统状态监控
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   //总处理日志行数
	Tps          float64 `json:"tps"`          //系统吞出量
	ReadChanLen  int     `json:"readChanLen"`  //read channel 长度
	WriteChanLen int     `json:"writeChanLen"` //write channel 长度
	RunTime      string  `json:"runTime"`      //运行总时间
	ErrNum       int     `json:"errNum"`       //错误数
}

var TypeMonitorChan = make(chan int, 200)

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

func (m *Monitor) Start(lp *proc.LogProcess) {

	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.Data.ErrNum += 1
			case TypeHandleLine:
				m.Data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-ticker.C
			m.TpsSlic = append(m.TpsSlic, m.Data.HandleLine)
			if len(m.TpsSlic) > 2 {
				m.TpsSlic = m.TpsSlic[1:]
			}

		}
	}()

	// 该服务的监控
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		m.Data.RunTime = time.Now().Sub(m.StartTime).String()
		m.Data.ReadChanLen = len(lp.RC)
		m.Data.WriteChanLen = len(lp.WC)

		if len(m.TpsSlic) > 2 {
			m.Data.Tps = float64(m.TpsSlic[1]-m.TpsSlic[0]) / 5
		}

		ret, _ := json.MarshalIndent(m.Data, "", "\t")
		io.WriteString(w, string(ret))

	})
	http.ListenAndServe(":9005", nil)
}
