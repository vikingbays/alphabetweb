// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"alphabet/log4go/message"
	"sync"
	"time"
)

/*
记录一个协程号的信息，包括：正在运行的协程号，序列号，该协程创建时间。
序列号信息，实际是用于记录一次调用链的过程，是一个轨迹信息。
*/
type Caller struct {
	gid          uint64    // 协程号
	serialnumber string    // 序列号
	createTime   time.Time // 当前gid创建的时间
}

func (c *Caller) GetGid() uint64 {
	return c.gid
}

func (c *Caller) GetSerialnumber() string {
	return c.serialnumber
}

func (c *Caller) GetCreateTime() time.Time {
	return c.createTime
}

// 用于测试的变量，如果是true，就会打印测试信息
var testing_chain bool = false

/*
调用链操作类
*/
type CallChain struct {
	callerMap        map[uint64]*Caller
	locker           *sync.RWMutex
	callchainChannel chan *Caller
	delGids          []uint64
	test_delGids     map[uint64]int
}

/*
添加一个协程信息，如果该协程已经被创建，那么就重新写入serialnumber信息
@param gid  协程号
@param serialnumber  序列号
*/
func (cc *CallChain) AddCaller(gid uint64, serialnumber string) {
	if serialnumber == "" {
		serialnumber = GetSerialnumber()
	}
	cc.locker.Lock()
	defer cc.locker.Unlock()
	if m1, ok := cc.callerMap[gid]; ok {
		m1.serialnumber = serialnumber
	} else {
		cc.callerMap[gid] = &Caller{gid: gid, serialnumber: serialnumber, createTime: time.Now()}
	}
}

/*
获取一个协程信息，根据gid
@param gid  协程号
*/
func (cc *CallChain) GetCaller(gid uint64) *Caller {
	cc.locker.RLock()
	defer cc.locker.RUnlock()
	return cc.callerMap[gid]
}

/*
获取当前的协程信息
*/
func (cc *CallChain) GetCurrentCaller() *Caller {
	gid := GetGID()
	cc.locker.RLock()
	defer cc.locker.RUnlock()
	return cc.callerMap[gid]
}

/*
删除一个协程信息，根据gid
@param gid  协程号
*/
func (cc *CallChain) DelCaller(gid uint64) {
	cc.locker.Lock()
	defer cc.locker.Unlock()
	delete(cc.callerMap, gid)
}

/*
批量删除一个协程信息，根据gid
@param gids  协程号 数组
*/
func (cc *CallChain) DelCaller_Batch(gids []uint64) {
	cc.locker.Lock()
	defer cc.locker.Unlock()
	for _, gid := range gids {
		delete(cc.callerMap, gid)
		if testing_chain {
			cc.test_delGids[gid] = 1
		}
	}
}

func (cc *CallChain) RecvCallchainChannel() {
	test_num := 0

	count := 1000
	for c := range cc.callchainChannel {
		cc.delGids = append(cc.delGids, c.gid)
		if len(cc.delGids) == count {
			delGids_tmp := make([]uint64, count, count)
			copy(delGids_tmp, cc.delGids)
			cc.delGids = make([]uint64, 0, 1)
			Goroutine(func() {
				if testing_chain {
					if len(delGids_tmp) != count {
						InfoLog("error:: len(delGids_tmp)=%d  , count=%d ", len(delGids_tmp), count)
					}

				}
				callchainGlobal.DelCaller_Batch(delGids_tmp)
			})

		}
		if testing_chain {
			test_num = test_num + 1
			if test_num%10000000 == 0 {
				InfoLog("info: testnum=%d", test_num)
			}
		}
	}
}

func (cc *CallChain) SendCallchainChannel(c *Caller) {
	cc.callchainChannel <- c
}

/*
检查协程信息，把已经无效的协程（已回收的），移除掉。
*/
func (cc *CallChain) Check() {
	defer func() { // 用于捕获panic异常，不影响整个服务运行
		if err := recover(); err != nil {
			ErrorLog(err)
		}
	}()

	map1 := GetAllGID()
	map2 := make(map[uint64]uint64)
	cc.locker.RLock()

	for k, _ := range cc.callerMap {
		if _, ok := map1[k]; !ok { // 如果map1中没有，说明该gid不存在了，标记需要删除
			map2[k] = k
		}
	}
	cc.locker.RUnlock()
	DebugLog(message.DEG_CALL_79010, cc.callerMap)
	cc.locker.Lock()
	for k, _ := range map2 {
		delete(cc.callerMap, k)
	}
	cc.locker.Unlock()
}
