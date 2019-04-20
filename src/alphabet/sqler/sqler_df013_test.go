// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package sqler

import (
	"alphabet/core/abtest"
	"testing"
)

// 重点测试 DoQueryMap 方法，获取的结果集数据格式是struct类型，例如： map[string]Case1_User1_df013
// 在此测试基础上增加 FilterRows 方法，增加对数据行的过滤处理，判断如果UsrId是奇数，那么就不记录（通过FilterRows返回值是false或者是0）。

//go test -v -test.run Test_sqler_df013
func Test_sqler_df013(t *testing.T) {
	dataFinder := test_init_sqler_df000(t)
	defer dataFinder.ReleaseConnection(0)

	test_sqler_df013_nest(dataFinder, t)

}

func test_sqler_df013_nest(dataFinder *DataFinder, t *testing.T) {

	count := 50
	test_sqler_df000_case1_initData(dataFinder, count, t)

	var minUserId, maxUserId int = 0, 100

	//paramUser := User1{MinUsrid: minUserId, MaxUsrid: maxUserId}
	paramUser := new(Case1_User1_df013)
	paramUser.MinUsrid = minUserId
	paramUser.MaxUsrid = maxUserId

	dataFinder.StartTrans()

	userMap := make(map[string]Case1_User1_df013)
	dataFinder.DoQueryMap("case1_query", paramUser, userMap, "Name", "")

	abtest.InfoTestLog(t, "len(userIdList)=%d", len(userMap))

	abtest.InfoTestLog(t, "(userMap)=%v", userMap)

	if len(userMap) != count/2 {
		abtest.ErrorTestLog(t, "查询结果（%d）!= 预设数据量(%d)", len(userMap), count/2)
	}

	dataFinder.CommitTrans()

}

type Case1_User1_df013 struct {
	Usrid    int
	Name     string
	Nanjing  bool
	Money    float64
	Hello    string
	MinUsrid int
	MaxUsrid int
}

func (usr *Case1_User1_df013) FilterRows() bool {

	usr.Name = "Apple_" + usr.Name
	if usr.Usrid%2 != 0 {
		return false
	} else {
		return true
	}
}
