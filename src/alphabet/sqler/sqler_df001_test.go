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

//go test -v -test.run Test_sqler_df001
// 重点测试 DoQueryList 方法，获取的结果集数据格式是基础类型，例如： []int
// 初始化 50条数据，调用DoQueryList，查看返回结果是不是 50
func Test_sqler_df001(t *testing.T) {
	dataFinder := test_init_sqler_df000(t)
	defer dataFinder.ReleaseConnection(0)

	test_sqler_df001_nest(dataFinder, t)

}

func test_sqler_df001_nest(dataFinder *DataFinder, t *testing.T) {

	count := 50
	test_sqler_df000_case1_initData(dataFinder, count, t)

	var minUserId, maxUserId int = 0, 100

	//paramUser := User1{MinUsrid: minUserId, MaxUsrid: maxUserId}
	paramUser := new(Case1_User1_df000)
	paramUser.MinUsrid = minUserId
	paramUser.MaxUsrid = maxUserId

	dataFinder.StartTrans()

	userIdList := make([]int32, 100, 300)
	dataFinder.DoQueryList("case1_query_userid", paramUser, &userIdList)

	abtest.InfoTestLog(t, "len(userIdList)=%d", len(userIdList))

	abtest.InfoTestLog(t, "(userIdList)=%v", userIdList)

	dataFinder.CommitTrans()

	if len(userIdList) != count {
		abtest.ErrorTestLog(t, "查询结果（%d）!= 预设数据量(%d)", len(userIdList), count)
	}

}
