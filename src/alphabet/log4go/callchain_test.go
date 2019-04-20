// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"sync"
	"testing"
	"time"
)

// 重点测试 发送 销毁消息和接收销毁消息是否一致？

//go test -v -test.run Test_Callchain
func Test_Callchain(t *testing.T) {
	initTestLog()

	testing_chain = true
	InitCallChain()

	count_gthread := 1000
	times_num := 100
	per_count := 100

	wg := new(sync.WaitGroup)
	wg.Add(1)

	wg2 := new(sync.WaitGroup)
	wg2.Add(count_gthread)

	InfoLog("start.....")
	for i := 0; i < count_gthread; i++ {
		sn_gthread := times_num * per_count * i
		go func() {
			wg.Wait()
			times := 0

			for true {
				if times == times_num {
					break
				}
				for j := 0; j < per_count; j++ {
					GetCallChain().SendCallchainChannel(&Caller{gid: uint64(sn_gthread + j + times*per_count)})
				}
				times = times + 1
				time.Sleep(time.Millisecond * 200)
			}
			wg2.Done()
		}()

	}

	time.Sleep(time.Second * 2)
	wg.Done()

	wg2.Wait()

	InfoLog("接收的数据：%d", len(callchainGlobal.test_delGids))
	//InfoLog("接收的数据：%v", callchainGlobal.test_delGids)
	InfoLog("This is End!")

	time.Sleep(time.Second * 100)
}
