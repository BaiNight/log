package read

import (
	"bufio"
	"fmt"
	"io"
	"logpro/defs"
	"os"
	"time"
)

type Reader interface {
	Read(chan<- []byte)
}

type FromFile struct {
	Path string
}

func (r *FromFile) Read(rc chan<- []byte) {
	f, err := os.Open(r.Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 将指针移动到文件末尾
	//f.Seek(0, 2)

	rd := bufio.NewReader(f)

	// 逐行读取文件
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("readbytes error: %s", err.Error()))
		}

		rc <- line

		defs.TypeMonitorChan <- defs.TypeHandleLine
	}
}
