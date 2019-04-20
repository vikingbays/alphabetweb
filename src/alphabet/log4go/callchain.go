// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"sync"
	"time"
)

var callchainGlobal *CallChain

func InitCallChain() {
	callchainGlobal = initCallChain_nest()

	Goroutine(func() {
		for true {
			time.Sleep(time.Second * 30)
			callchainGlobal.Check()
		}

	})

	Goroutine(func() {
		callchainGlobal.RecvCallchainChannel()
	})

}

func initCallChain_nest() *CallChain {
	callChain := new(CallChain)
	callChain.callerMap = make(map[uint64]*Caller)
	callChain.locker = new(sync.RWMutex)
	callChain.callchainChannel = make(chan *Caller, 100)
	callChain.delGids = make([]uint64, 0, 1)
	if testing_chain {
		callChain.test_delGids = make(map[uint64]int)
	}
	return callChain
}

func GetCallChain() *CallChain {
	return callchainGlobal
}

func NewCaller(gid uint64, serialnumber string, createTime time.Time) *Caller {
	return &Caller{gid: gid, serialnumber: serialnumber, createTime: createTime}
}

/*
获取当前协程的唯一序列号
*/
func GetUSN() string {
	cc := GetCallChain()
	if cc == nil {
		return ""
	} else {
		c := cc.GetCurrentCaller()
		if c != nil {
			return c.GetSerialnumber()
		}
	}
	return ""
}
