// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/abtest"
	"alphabet/env"
	"sync"
	"testing"
	"time"
)

// 测试 一个 http2连接，支撑 200个请求的并发场景，效果非常理想. 建议设置100以下
// 调节并发数             num_parallel := 200
// 设置验证的返回结果长度   len_resp := 8552
//

// go test -v -test.run Test_conn_http2
func Test_conn_http2(t *testing.T) {
	abtest.InitTestEnv()
	abtest.InitTestLog()
	env.Switch_CallChain = false

	//sc := &SimpleRpcClient{Protocol: "rpc_tcp_ssl", Ip: "127.0.0.1", Port: 11443}
	sc := &SimpleRpcClient{Protocol: "rpc_tcp_ssl", Ip: "localhost", Port: 10766}
	err := sc.Connect()
	if err != nil {
		t.Log(err)
	}

	defer sc.Close()

	num_parallel := 200
	len_resp := 8552

	wg := new(sync.WaitGroup)
	wg2 := new(sync.WaitGroup)
	wg.Add(1)
	wg2.Add(num_parallel)
	t1 := time.Now()
	for i := 0; i < num_parallel; i++ {
		go func() {
			wg.Wait()
			d1, err2 := sc.DoPostDatas("/octopus_service/db/query2", "min=0&max=100")
			if err2 != nil {
				t.Error(err2)
			} else {
				if len_resp != len(d1) {
					t.Error(len(d1))
				}
				t.Log(len(d1))
			}
			wg2.Done()
		}()
	}
	wg.Done()
	wg2.Wait()
	t2 := time.Now()
	time.Sleep(time.Second * 1)
	t.Logf("耗时：%d ms", (t2.UnixNano()-t1.UnixNano())/1000/1000)

	//time.Sleep(time.Second * 100)
}

// go test -v -test.run Test_conn_http2
func Test_conn_head(t *testing.T) {
	abtest.InitTestEnv()
	abtest.InitTestLog()
	env.Switch_CallChain = false

	//sc := &SimpleRpcClient{Protocol: "rpc_tcp_ssl", Ip: "127.0.0.1", Port: 11443}
	sc := &SimpleRpcClient{Protocol: "rpc_tcp_ssl", Ip: "localhost", Port: 10766}
	err := sc.Connect()
	if err != nil {
		t.Log(err)
	}

	defer sc.Close()

	d1, err2 := sc.head_nest("/", true)

	if err2 != nil {
		t.Error(err2)
	} else {
		t.Log(len(d1))
	}

	d1, err2 = sc.DoPostDatas("/octopus_service/db/query2", "min=0&max=100")

	if err2 != nil {
		t.Error(err2)
	} else {
		t.Log(len(d1))
	}

	t.Log(sc.client)

}
