// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/abtest"
	"alphabet/env"
	"bytes"
	"sync"
	"testing"
	"time"
)

// 测试 一个 http1.1连接，支撑 10个请求的并发场景，没法并行，直接报错了

// go test -v -test.run Test_conn_http11
func Test_conn_http11(t *testing.T) {
	abtest.InitTestEnv()
	abtest.InitTestLog()
	env.Switch_CallChain = false

	//sc := &SimpleRpcClient{Protocol: "rpc_tcp_ssl", Ip: "127.0.0.1", Port: 11443}
	sc := &SimpleRpcClient{Protocol: "rpc_tcp", Ip: "localhost", Port: 10777}
	err := sc.Connect()
	if err != nil {
		t.Log(err)
	}

	defer sc.Close()

	num_parallel := 10
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
}

// go test -v -test.run Test_bytes
func Test_bytes(t *testing.T) {
	p := make([]byte, 6)
	readP(p)
	t.Log(p)
}

func readP(p []byte) {
	formFieldBuf1 := &bytes.Buffer{}
	formFieldBuf1.Write([]byte{'A', 'B', 'C', 'D'})

	formFieldBuf2 := &bytes.Buffer{}
	formFieldBuf2.Write([]byte{'E', 'F', 'G', 'H'})

	n, _ := formFieldBuf1.Read(p)
	if n < len(p) {
		formFieldBuf2.Read(p[n:])
	}
}
