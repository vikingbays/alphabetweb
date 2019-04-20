// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"testing"
	"time"
)

//测试添加  %U，用来记录唯一序列号的功能
//go test -v -test.run Test_USN
func Test_USN(t *testing.T) {
	initTestLog()
	time.Sleep(time.Second * 1)
	InfoLog("aaaaaaaa")

	time.Sleep(time.Second * 2)
}
