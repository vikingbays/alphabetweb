// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package etcd

import (
	"alphabet/env"
	"alphabet/log4go"
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	dialTimeout = 5 * time.Second
	endpoints   = []string{"127.0.0.1:2379"}
	username    = "user1"
	passwd      = "user1pw"
	rootPath    = "/ssfoo/"
)

func testInit(t *testing.T) {
	projectPath := os.Getenv("GOLANG_PROJECT_TESTUNIT")

	t.Log("projectPath:" + projectPath)

	env.InitWithoutLoad(projectPath)

	logconfigFile := env.Env_Project_Resource + "/alphabet/log4go/logconfig_test.toml"
	t.Log("logconfigFile:" + logconfigFile)

	log4go.InitLogger(logconfigFile, env.Env_Project_Root)
}

//go test -v -bench="etcdconn_test.go"
// export GOLANG_PROJECT_TESTUNIT="/Users/xxx/devs/golang/projects/AlphabetwebProject/alphabetweb"

// 测试连接
func Test_EtcdConn(t *testing.T) {
	testInit(t)
	etcdConnection := EtcdConnection{}

	// 测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		etcdConnection.Put(rootPath+"sample1_key", "data1_value1")

		v11, err11 := etcdConnection.Get(rootPath + "sample1_key")
		if v11 != "data1_value1" || err11 != nil {
			t.Errorf("err: [0]Get value after Put() , result: %s   ,    error_info: %s ", v11, err11)
		}
		etcdConnection.Del(rootPath + "sample1_key")
		v12, err12 := etcdConnection.Get(rootPath + "sample1_key")
		if v12 == "data1_value1" || err11 != nil {
			t.Errorf("err: [1]Get value after Del() , result: %s   ,    error_info: %s ", v12, err12)
		}
		etcdConnection.Close()
	}

	// 测试 连接未打开时，  put
	{
		err11 := etcdConnection.Put(rootPath+"sample1_key", "data1_value1")
		if err11 == nil {
			t.Errorf("err: [2]Put() without Connecting ,    error_info: %s ", err11)
		}
	}

	// 测试 连接未打开时，  put
	{
		err11 := etcdConnection.Del(rootPath + "sample1_key")
		if err11 == nil {
			t.Errorf("err: [3]Del() without Connecting ,    error_info: %s ", err11)
		}
	}

	// 测试 连接未打开时，  get
	{
		v11, err11 := etcdConnection.Get(rootPath + "sample1_key")
		if err11 == nil {
			t.Errorf("err: [4]Get() without Connecting  , result: %s   ,    error_info: %s ", v11, err11)
		}
	}

	// 重新连接，再测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		etcdConnection.Put(rootPath+"sample1_key", "data1_value1")

		v11, err11 := etcdConnection.Get(rootPath + "sample1_key")
		if v11 != "data1_value1" || err11 != nil {
			t.Errorf("err: [0]Get value after Put() , result: %s   ,    error_info: %s ", v11, err11)
		}
		etcdConnection.Del(rootPath + "sample1_key")
		v12, err12 := etcdConnection.Get(rootPath + "sample1_key")
		if v12 == "data1_value1" || err11 != nil {
			t.Errorf("err: [1]Get value after Del() , result: %s   ,    error_info: %s ", v12, err12)
		}
		etcdConnection.Close()
	}

}

// 测试事务情况
func Test_EtcdTrans(t *testing.T) {

	testInit(t)
	etcdConnection := EtcdConnection{}

	etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)

	etcdConnection.Del(rootPath + "sample2_key001")
	etcdConnection.Del(rootPath + "sample2_key002")
	etcdConnection.Del(rootPath + "sample2_key003")
	etcdConnection.Del(rootPath + "sample2_key004")
	etcdConnection.Del(rootPath + "sample2_key005")

	etcdConnection.StartTrans()
	etcdConnection.Put(rootPath+"sample2_key001", "data1_value1_for_key001")
	etcdConnection.Put(rootPath+"sample2_key002", "data1_value2_for_key002")
	etcdConnection.Put(rootPath+"sample2_key003", "data1_value3_for_key003")
	etcdConnection.Put(rootPath+"sample2_key004", "data1_value4_for_key004")
	etcdConnection.Put(rootPath+"sample2_key005", "data1_value5_for_key005")
	etcdConnection.Del(rootPath + "sample2_key003")

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample2_key001")
		if v11 != "" || err11 != nil {
			t.Errorf("err: [0]Get value after Put() ,but trans is not committed , result: %s   ,    error_info: %s ", v11, err11)
		}
	}
	etcdConnection.CommitTrans()

	v11, err11 := etcdConnection.Get(rootPath + "sample2_key001")
	if v11 != "data1_value1_for_key001" || err11 != nil {
		t.Errorf("err: [1]Get value after commit trans , result: %s   ,    error_info: %s ", v11, err11)
	} else {
		t.Log(v11)
	}

	v12, err12 := etcdConnection.Get(rootPath + "sample2_key002")
	if v12 != "data1_value2_for_key002" || err12 != nil {
		t.Errorf("err: [2]Get value after commit trans , result: %s   ,    error_info: %s ", v12, err12)
	}

	v13, err13 := etcdConnection.Get(rootPath + "sample2_key003")
	if v13 != "" || err13 != nil {
		t.Errorf("err: [3]Get value after commit trans , result: %s   ,    error_info: %s ", v13, err13)
	}

	v14, err14 := etcdConnection.Get(rootPath + "sample2_key004")
	if v14 != "data1_value4_for_key004" || err14 != nil {
		t.Errorf("err: [4]Get value after commit trans , result: %s   ,    error_info: %s ", v14, err14)
	}

	v15, err15 := etcdConnection.Get(rootPath + "sample2_key005")
	if v15 != "data1_value5_for_key005" || err15 != nil {
		t.Errorf("err: [5]Get value after commit trans , result: %s   ,    error_info: %s ", v15, err15)
	}

	etcdConnection.Close()

}

// 测试TTL
func Test_EtcdTTL(t *testing.T) {

	testInit(t)
	etcdConnection := EtcdConnection{}

	etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
	etcdConnection.Del(rootPath + "sample3_key001")
	etcdConnection.Del(rootPath + "sample3_key002")
	etcdConnection.Del(rootPath + "sample3_key003")

	var timeout int64 = 4
	respTTL, err := etcdConnection.TTL_SetGrant(timeout)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	} else {
		err0 := etcdConnection.TTL_Put(rootPath+"sample3_key001", "value1_TTL", respTTL)
		if err0 != nil {
			t.Errorf("err: %s", err0.Error())
		}
	}

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample3_key001")
		if v11 != "value1_TTL" || err11 != nil {
			t.Errorf("err: [11]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}
	}

	time.Sleep(time.Duration((timeout - 1) * 1000 * 1000 * 1000)) // (timeout - 1) 秒

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample3_key001")
		if v11 != "value1_TTL" || err11 != nil {
			t.Errorf("err: [12]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}

		d1, d2, err := etcdConnection.TTL_GetTimeToLive(respTTL)
		t.Logf(" Total Time: %d, Live Time： %d,err: %v", d1, d2, err)

	}

	time.Sleep(2 * time.Second)

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample3_key001")
		if v11 != "" || err11 != nil {
			t.Errorf("err: [13]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}
		d1, d2, err := etcdConnection.TTL_GetTimeToLive(respTTL)
		t.Logf(" Total Time: %d, Live Time： %d,err: %v", d1, d2, err)

	}

	etcdConnection.Close()

}

// 测试TTL 带事务
func Test_EtcdTTLAndTrans(t *testing.T) {

	testInit(t)
	etcdConnection := EtcdConnection{}

	etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
	etcdConnection.Del(rootPath + "sample4_key001")
	etcdConnection.Del(rootPath + "sample4_key002")
	etcdConnection.Del(rootPath + "sample4_key003")

	etcdConnection.StartTrans()

	var timeout int64 = 4
	respTTL, err := etcdConnection.TTL_SetGrant(timeout)
	if err != nil {
		t.Errorf("err: %s", err.Error())
	} else {
		err0 := etcdConnection.TTL_Put(rootPath+"sample4_key001", "value1_TTL", respTTL)
		if err0 != nil {
			t.Errorf("err: %s", err0.Error())
		}
	}

	etcdConnection.CommitTrans()

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample4_key001")
		if v11 != "value1_TTL" || err11 != nil {
			t.Errorf("err: [11]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}
	}

	time.Sleep(time.Duration((timeout - 1) * 1000 * 1000 * 1000)) // (timeout - 1) 秒

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample4_key001")
		if v11 != "value1_TTL" || err11 != nil {
			t.Errorf("err: [12]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}

		d1, d2, err := etcdConnection.TTL_GetTimeToLive(respTTL)
		t.Logf(" Total Time: %d, Live Time： %d,err: %v", d1, d2, err)

	}

	time.Sleep(2 * time.Second)

	{
		v11, err11 := etcdConnection.Get(rootPath + "sample4_key001")
		if v11 != "" || err11 != nil {
			t.Errorf("err: [13]Get value after TTL , result: %s   ,    error_info: %s ", v11, err11)
		}
		d1, d2, err := etcdConnection.TTL_GetTimeToLive(respTTL)
		t.Logf(" Total Time: %d, Live Time： %d,err: %v", d1, d2, err)

	}

	etcdConnection.Close()

}

var chan1_Watch chan opDatas = make(chan opDatas, 1)
var chan2_Watch chan opDatas = make(chan opDatas, 1)

// 测试watch
func Test_EtcdWatch(t *testing.T) {

	testInit(t)
	etcdConnection0 := EtcdConnection{}

	etcdConnection0.CreateConnection(endpoints, dialTimeout, username, passwd)
	defer etcdConnection0.Close()
	go func() {
		ctx, _ := context.WithCancel(context.Background())
		err1 := etcdConnection0.WatchWithPrefix(ctx, rootPath+"sample5_key001", putEventFunc, deleteEventFunc, nil)
		if err1 != nil {
			t.Errorf("err: %s", err1.Error())
		}
	}()
	go func() {
		ctx, _ := context.WithCancel(context.Background())
		err2 := etcdConnection0.WatchWithPrefix(ctx, rootPath+"sample5_key002", putEventFunc, deleteEventFunc, nil)
		if err2 != nil {
			t.Errorf("err: %s", err2.Error())
		}
	}()

	opDatasArray1 := make([]opDatas, 0, 4)
	opDatasArray2 := make([]opDatas, 0, 4)

	go func() {

		for {
			select {
			case v1, _ := <-chan1_Watch:
				opDatasArray1 = append(opDatasArray1, v1)
			case v2, _ := <-chan2_Watch:
				opDatasArray2 = append(opDatasArray2, v2)
			}
		}
	}()

	go func() {
		etcdConnection := EtcdConnection{}
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		defer etcdConnection.Close()
		etcdConnection.Put(rootPath+"sample5_key001", "value1_0000")
		for i := 0; i < 10; i++ {
			if i%3 == 0 || i%3 == 1 {
				time.Sleep(1 * time.Second)
				etcdConnection.StartTrans()
				etcdConnection.Put(rootPath+"sample5_key001", "value1_000"+string(i))
				etcdConnection.CommitTrans()
			} else {
				time.Sleep(1 * time.Second)
				etcdConnection.Del(rootPath + "sample5_key001")
			}

		}
	}()

	go func() {
		etcdConnection := EtcdConnection{}
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		defer etcdConnection.Close()

		var timeout int64 = 4
		respTTL, err := etcdConnection.TTL_SetGrant(timeout)
		if err != nil {
			t.Errorf("err: %s", err.Error())
		} else {
			err0 := etcdConnection.TTL_Put(rootPath+"sample5_key002", "value2_TTL_0000", respTTL)
			if err0 != nil {
				t.Errorf("err: %s", err0.Error())
			}
		}

		for i := 0; i < 4; i++ {
			time.Sleep(1 * time.Second)
			etcdConnection.TTL_SetKeepAliveOnce(respTTL)
			etcdConnection.StartTrans()
			etcdConnection.TTL_Put(rootPath+"sample5_key002", fmt.Sprintf("value2_TTL_000%d", (i+1)), respTTL)
			etcdConnection.CommitTrans()
		}
	}()

	time.Sleep(15 * time.Second)

	for _, data := range opDatasArray1 {
		t.Logf("[opDatasArray1]  oper: %d , key : %s , value : %s \n", data.oper, data.key, data.value)
	}

	for _, data := range opDatasArray2 {
		t.Logf("[opDatasArray2]  oper: %d , key : %s , value : %s \n", data.oper, data.key, data.value)
	}

	if len(opDatasArray1) != 8 {
		t.Error("err : opDatasArray1")
	}

	if len(opDatasArray2) != 6 {
		t.Error("err : opDatasArray2")
	}

}

func putEventFunc(key string, value string) {
	fmt.Printf("Put:  {key: %s , value :%s} \n", key, value)
	if key == rootPath+"sample5_key001" {
		chan1_Watch <- opDatas{oper: 0, key: key, value: value}
	} else {
		chan2_Watch <- opDatas{oper: 0, key: key, value: value}
	}

}

func deleteEventFunc(key string, value string) {
	fmt.Printf("Delete:  {key: %s , value :%s} \n", key, value)
	if key == rootPath+"sample5_key001" {
		chan1_Watch <- opDatas{oper: 1, key: key, value: value}
	} else {
		chan2_Watch <- opDatas{oper: 1, key: key, value: value}
	}
}

type opDatas struct {
	oper  int
	key   string
	value string
}
