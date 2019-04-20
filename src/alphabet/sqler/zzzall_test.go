// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/core/abtest"
	"testing"
	"time"
)

type TestCase struct {
	TestName string             //测试名称
	TestFunc func(t *testing.T) // 测试方法
	Testable bool               // 是否测试
	Result   bool               // 测试是否通过
}

func initTestSqler() []*TestCase {
	testCases := []*TestCase{
		&TestCase{TestName: "Test_sqler_df001", TestFunc: Test_sqler_df001, Testable: true},
		&TestCase{TestName: "Test_sqler_df002", TestFunc: Test_sqler_df002, Testable: true},
		&TestCase{TestName: "Test_sqler_df003", TestFunc: Test_sqler_df003, Testable: true},
		&TestCase{TestName: "Test_sqler_df004", TestFunc: Test_sqler_df004, Testable: true},
		&TestCase{TestName: "Test_sqler_df011", TestFunc: Test_sqler_df011, Testable: true},
		&TestCase{TestName: "Test_sqler_df012", TestFunc: Test_sqler_df012, Testable: true},
		&TestCase{TestName: "Test_sqler_df013", TestFunc: Test_sqler_df013, Testable: true},
		&TestCase{TestName: "Test_sqler_df091", TestFunc: Test_sqler_df091, Testable: true},
		&TestCase{TestName: "Test_connectionpool001", TestFunc: Test_connectionpool001, Testable: true},
		&TestCase{TestName: "Test_connectionpool002", TestFunc: Test_connectionpool002, Testable: true},
	}
	return testCases
}

//go test -v -test.run Test_ALL
// 测试sqler下的所有测试用例
func Test_ALL(t *testing.T) {
	testCases := initTestSqler()
	errIndexs := make([]int, 0, 1)
	for i, tc := range testCases {
		tc.Result = t.Run(tc.TestName, tc.TestFunc)
		if !tc.Result {
			errIndexs = append(errIndexs, i)
		}
	}
	abtest.InfoTestLog(t, "Alphabetweb.sqler.Case , success: %d , failed: %d .", (len(testCases) - len(errIndexs)), len(errIndexs))
	for i, index := range errIndexs {
		abtest.InfoTestLog(t, "Alphabetweb.sqler.Case , error[%d]: %s", i, testCases[index].TestName)
	}
	time.Sleep(time.Second * 1)
}
