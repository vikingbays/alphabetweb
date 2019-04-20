// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)] // 如果是false，表示只是当前的协程号；如果true，表示所有有效的协程号
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func GetAllGID() map[uint64]string {
	gidAndStatusMap := make(map[uint64]string)

	buf := make([]byte, 1<<20)
	runtime.Stack(buf, true)
	for true {
		index1 := bytes.Index(buf, []byte("goroutine "))
		if index1 == -1 {
			break
		}
		if len(buf) > index1+10 {
			buf = buf[index1+10:]
			index2 := bytes.IndexByte(buf, ' ')
			if len(buf) > index2 {
				n_gid, _ := strconv.ParseUint(string(buf[:index2]), 10, 64)

				n_status := ""
				if len(buf) > index2+2 {
					buf = buf[index2+2:]
					index3 := bytes.IndexByte(buf, ']')
					if index3 != -1 {
						n_status = string(buf[:index3])
					}
				} else {
					break
				}

				gidAndStatusMap[n_gid] = n_status
			} else {
				break
			}
		} else {
			break
		}
	}
	return gidAndStatusMap
}

/*
获取一个新的序列号
*/
func GetSerialnumber() string {
	var ticketFromRand string = strconv.FormatInt(time.Now().UnixNano(), 16)
	for i := 0; i < 20; i++ {
		numFromChar := 65 + rand.Intn(25)
		ticketFromRand = fmt.Sprintf("%s%s", ticketFromRand, string(rune(numFromChar)))
	}
	return ticketFromRand
}

/*
Goroutine 创建
*/
func Goroutine(func_sample func()) {
	gid_parent := GetGID()
	caller := GetCallChain().GetCaller(gid_parent)
	serialnumber := ""
	if caller != nil {
		serialnumber = caller.GetSerialnumber()
	}
	go func(sn string) {
		gid_nest := GetGID()
		callchainGlobal.AddCaller(gid_nest, serialnumber)
		func_sample()
	}(serialnumber)
}
