// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/abtest"
	"alphabet/env"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// micro-service配置
// type="etcd"
// endpoints = ["127.0.0.1:2379"]
// username  = "serv001"
// password  = "123456"
// timeout = 2     # 2秒
// root = "awroot_serv001"
//
// 初始化
// export  ETCDCTL_API=3
// ./etcd3/etcdctl --endpoints="localhost:2379"  user add  root
// 设置密码【rootpw】
// ./etcd3/etcdctl --endpoints="localhost:2379" auth enable
// ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw role add role1
// ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  role grant-permission role1 --prefix=true  readwrite  awroot_serv001
// ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  user add serv001
// 设置密码【123456】
// ./etcd3/etcdctl --endpoints="localhost:2379" --user root:rootpw  user grant-role  serv001  role1

// 验证：
//./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   get  --prefix   awroot_serv001

// ------

// export  ETCDCTL_API=3
//./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   get  --prefix   awroot_serv001

//./etcd3/etcdctl --endpoints="localhost:2379"  --user serv001:123456  lease grant 20
//./etcd3/etcdctl --endpoints="localhost:2379"  --user serv001:123456  put "awroot_serv001.groups.g1.addrs.[127.0.0.1_9000]" ttl --lease=694d64212e07b997

//go test -v -bench="configforservcli_test.go"
// go test -v -test.run Test_initServerAll
// export GOLANG_PROJECT_TESTUNIT="/Users/xxx/devs/golang/projects/AlphabetwebProject/alphabetweb"

/*
测试 客户端监听。

首先启动客户端监听测试：  go test -v -test.run Test_initClientAll
然后启动server端注册：   go test -v -test.run Test_initServerAll

1、 测试新增一个group
      awroot_serv001.groups.g1
		（启动server 注册即可：go test -v -test.run Test_initServerAll ）
2、 测试监听一个addr
      awroot_serv001.groups.g1.addrs.[127.0.0.1_9000]
3、 测试监听一个ticket
      awroot_serv001.groups.g1.addrs.[127.0.0.1_9000].ticket
      awroot_serv001.groups.g1.addrs.[127.0.0.1_9000].ticketforlast

./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   put  "awroot_serv001.groups.g1.addrs.[127.0.0.1_9000].ticketforlast" "aaaaa"
./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   put  "awroot_serv001.groups.g1.addrs.[127.0.0.1_9000].ticket" "bbbbbb"


4、 测试timestamp 更新
     awroot_serv001.groups.g1.timestamp
./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   put  "awroot_serv001.groups.g1.timestamp"  "20180629171627621847000"

5、 测试timestamp 删除
     awroot_serv001.groups.g1.timestamp
./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   del  "awroot_serv001.groups.g1.timestamp"
./etcd3/etcdctl --endpoints="localhost:2379" --user serv001:123456   del  --prefix "awroot_serv001.groups.g1"

*/

func Test_initServerAll(t *testing.T) {
	abtest.InitTestEnv()
	abtest.InitTestLog()

	serverGlobalMap = make(map[string]*Server_Info)
	serverGlobalMap["alphabet"] = new(Server_Info)
	pathconfig := env.Env_Project_Resource + "/" + "alphabet/service/sample_server_test.toml"
	initServerAll(pathconfig, "sample_server_test.toml", "alphabet")
}

func Test_initClientAll(t *testing.T) {
	abtest.InitTestEnv()
	abtest.InitTestLog()

	rootAppname := "alphabet"

	clientsGlobalMap = make(map[string]*Clients_Info)

	clientsGlobalMap[rootAppname] = new(Clients_Info)
	clientsGlobalMap[rootAppname].rootAppname = rootAppname
	clientsGlobalMap[rootAppname].configMap = make(map[string]Client_Config_List_Toml)
	clientsGlobalMap[rootAppname].managerMap = make(map[string]IServiceManager)
	clientsGlobalMap[rootAppname].dataMap = make(map[string]*ServiceManagerGroup)
	clientsGlobalMap[rootAppname].groupidToFilenameMap = make(map[string]string)

	pathconfig := env.Env_Project_Resource + "/" + "alphabet/service/sample_client_test.toml"
	initClientAll(pathconfig, "sample_client_test.toml", rootAppname)
	t.Log(clientsGlobalMap[rootAppname].dataMap["g1"])

	time.Sleep(1000 * time.Second)
}

func Test_InitForServerConfig(t *testing.T) {

	abtest.InitTestEnv()
	abtest.InitTestLog()

	pathconfig := env.Env_Project_Resource + "/" + "alphabet/service/sample_server_test.toml"

	serverConfigListToml := initForServerConfig(pathconfig, "alphabet")

	t.Log(serverConfigListToml)

	if len(serverConfigListToml.Server) == 0 {
		t.Error("error :  serverConfigListToml  is nil")
	} else if serverConfigListToml.Server[0].Register[0].Type != "etcd" {
		t.Error("error :  serverConfigListToml.server.register.type != 'etcd'  ")
	}
}

func Test_InitForClientConfig(t *testing.T) {

	abtest.InitTestEnv()
	abtest.InitTestLog()

	pathconfig := env.Env_Project_Resource + "/" + "alphabet/service/sample_client_test.toml"

	clientConfigListToml := initForClientConfig(pathconfig)

	t.Log(clientConfigListToml)

	if len(clientConfigListToml.Client) == 0 {
		t.Error("error :  clientConfigListToml  is nil")
	} else if len(clientConfigListToml.Client[0].GroupIds) == 0 {
		t.Error("error : len( clientConfigListToml.GroupIds ) == 0  ")
	}

	ticketFromRand := ""
	for i := 0; i < 20; i++ {
		numFromChar := 126 - rand.Intn(93)
		ticketFromRand = fmt.Sprintf("%s%s", ticketFromRand, string(rune(numFromChar)))
	}

	fmt.Printf("ticketFromRand: %s\n", ticketFromRand)

}

func Test_Copy(t *testing.T) {
	a1 := make([]string, 0, 1)
	a1 = append(a1, "jack")
	a1 = append(a1, "nacy")
	a1 = append(a1, "mickel")
	a1 = append(a1, "black")
	t.Log(a1)
	a2 := make([]string, len(a1), len(a1))
	copy(a2, a1)
	t.Log(a2)
	a1 = make([]string, 0, 1)
	t.Log(a1)
	t.Log(a2)
}

func Test_fileWalk(t *testing.T) {
	path := "/Users/vikingbays/golang/AlphabetwebProject/out/src/alphabetsample"

	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return err
		}
		t.Log(path)
		return nil
	})

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			t.Log(">>>>>>>>" + file.Name())
		}
	}
}
