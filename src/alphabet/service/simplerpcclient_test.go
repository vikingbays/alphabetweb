// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

type RespBody struct {
	Err  string     `json:"err"`
	Json []UserInfo `json:"json"`
}

type UserInfo struct {
	Usrid   int     `json:"Usrid"`
	Name    string  `json:"Name"`
	Nanjing bool    `json:"Nanjing"`
	Money   float64 `json:"Money"`
}

// go test -v -test.run Test_rpc_sample1
func Test_rpc_sample1(t *testing.T) {
	num := 1
	for i := 0; i < num; i++ {
		go test_rpc_sample1_001(t)
		//	go test_rpc_sample1_001(t)
		//	go test_rpc_sample1_001(t)
		//	go test_rpc_sample1_001(t)
		//	go test_rpc_sample1_001(t)
	}
	time.Sleep(10 * time.Second)
}

func test_rpc_sample1_001(t *testing.T) {
	t.Log("开始测试！")

	sc := &SimpleRpcClient{Protocol: "rpc_tcp", Ip: "127.0.0.1", Port: 7777}
	err := sc.Connect()
	if err != nil {
		t.Log(err)
	}

	defer sc.Close()

	respBody := RespBody{}
	sc.DoJson("/web2/restful/jsoninfo/0/1", "", &respBody)

	t.Log(respBody)

	_, errUpload := sc.PreUpload("/web2/upload/uploadfile").
		AddFile("path", "~/Downloads/postgresql-9.4.1209.jar", "pgsql20009.jar").
		AddParam("alias", "pgsql2.jar").
		AddParam("author", "jack2").
		AddParam("name", "n1").
		AddParam("name", "n2").DoUpload()

	if errUpload != nil {
		t.Log(errUpload)
	}

	sc.DoDownload("/web2/download/do_download", "filepath=~/Downloads/postgresql-9.4.1209.jar&aliasname=pg.jar", "~/devs/golang/projects/AlphabetwebProject/alphabetsample/upload/upload/tu.jar")

	reader, errStream := sc.DoStreamDown("/web2/db/query", strings.NewReader("min=0&max=100"))
	if errStream != nil {
		t.Log(errStream)
	}
	bytes, errBytes := ioutil.ReadAll(reader)
	if errBytes != nil {
		t.Log(errBytes)
	}
	defer reader.Close()
	t.Log(string(bytes))

}
