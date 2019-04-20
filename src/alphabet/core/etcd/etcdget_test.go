// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package etcd

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

//go test -v -bench="etcdget_test.go"
// go test -v -test.run Test_EtcdGet
// export GOLANG_PROJECT_TESTUNIT="/Users/xxx/devs/golang/projects/AlphabetwebProject/alphabetweb"

// 测试连接
func Test_EtcdGet(t *testing.T) {
	testInit(t)
	etcdConnection := EtcdConnection{}

	// 测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)

		etcdConnection.Put(rootPath+"sample1_key1", "data1_value1")
		etcdConnection.Put(rootPath+"sample1_key2", "data1_value2")
		etcdConnection.Put(rootPath+"sample1_key3", "data1_value3")
		etcdConnection.Put(rootPath+"sample1_key4", "data1_value4")

		kvMaps, _ := etcdConnection.GetWithPrefix(rootPath + "sample1_key")

		for k, v := range kvMaps {
			fmt.Printf("k=%s , v=%s \n", k, v)
		}

		v1, _ := etcdConnection.Get(rootPath + "sample1_key")
		if v1 == "" {
			fmt.Println("v1 is ''")
		}

		etcdConnection.Close()
	}

}

func Test_EtcdDelTran(t *testing.T) {
	testInit(t)
	etcdConnection := EtcdConnection{}

	// 测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)

		{
			kvMaps, _ := etcdConnection.GetWithPrefix(rootPath + "sample1_key")

			for k, v := range kvMaps {
				fmt.Printf("k=%s , v=%s \n", k, v)
			}

			v1, _ := etcdConnection.Get(rootPath + "sample1_key")
			if v1 == "" {
				fmt.Println("v1 is ''")
			}
		}

		etcdConnection.StartTrans()
		etcdConnection.DelWithPrefix(rootPath + "sample1_key")
		etcdConnection.CommitTrans()

		{
			kvMaps, _ := etcdConnection.GetWithPrefix(rootPath + "sample1_key")

			for k, v := range kvMaps {
				fmt.Printf("k=%s , v=%s \n", k, v)
			}

			v1, _ := etcdConnection.Get(rootPath + "sample1_key")
			if v1 == "" {
				fmt.Println("v1 is ''")
			}
		}
		etcdConnection.Close()
	}

}

func Test_EtcdDel(t *testing.T) {
	testInit(t)
	etcdConnection := EtcdConnection{}

	// 测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		etcdConnection.DelWithPrefix(rootPath + "sample1_key")

		kvMaps, _ := etcdConnection.GetWithPrefix(rootPath + "sample1_key")

		for k, v := range kvMaps {
			fmt.Printf("k=%s , v=%s \n", k, v)
		}

		v1, _ := etcdConnection.Get(rootPath + "sample1_key")
		if v1 == "" {
			fmt.Println("v1 is ''")
		}

		etcdConnection.Close()
	}

}

func Test_EtcdGetNext(t *testing.T) {
	testInit(t)
	etcdConnection := EtcdConnection{}

	// 测试 正常的  put，get,del ,get
	{
		etcdConnection.CreateConnection(endpoints, dialTimeout, username, passwd)
		etcdConnection.Del(rootPath + "sample1_key1")
		//etcdConnection.Put(rootPath+"sample1_key1", "data1_value1")
		etcdConnection.Put(rootPath+"sample1_key2", "data1_value2")
		etcdConnection.Put(rootPath+"sample1_key3", "data1_value3")
		etcdConnection.Put(rootPath+"sample1_key4", "data1_value4")

		etcdConnection.Put(rootPath+"sample1_key1.sss", "data1_value1")

		v1, _ := etcdConnection.Get(rootPath + "sample1_key1.*")

		fmt.Printf("v1 is %s \n", v1)

		etcdConnection.Close()
	}
}

func Test_S001(t *testing.T) {
	s1 := "192.168.1.22_9000,fe80::cbd:8cd5:ec3a:160b%en0_8000"
	array := strings.Split(s1, ",")
	for _, a := range array {
		fmt.Println(a)

		array2 := strings.Split(a, "_")

		fmt.Println(array2[0])
		fmt.Println(array2[1])
	}
}

// 测试watch
func Test_EtcdWatch001(t *testing.T) {

	testInit(t)
	etcdConnection0 := EtcdConnection{}

	etcdConnection0.CreateConnection(endpoints, dialTimeout, username, passwd)
	defer etcdConnection0.Close()
	ctx, _ := context.WithCancel(context.Background())

	go func() {
		err1 := etcdConnection0.Watch(ctx, rootPath+"sample5_key001", putEventFunc0, deleteEventFunc0, errFunc0)
		if err1 != nil {
			t.Errorf("err: %s", err1.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	etcdConnection0.Put(rootPath+"sample5_key001", "33333")

	etcdConnection0.Put(rootPath+"sample5_key001", "44444")

	time.Sleep(1 * time.Second)

	//etcdConnection0.Close()

	fmt.Println(etcdConnection0.IsOpened())

	time.Sleep(30 * time.Second)

	fmt.Println(etcdConnection0.IsOpened())

	etcdConnection0.Put(rootPath+"sample5_key001", "55555")

	time.Sleep(10 * time.Second)

}

func putEventFunc0(key string, value string) {
	fmt.Printf("Put:  {key: %s , value :%s} \n", key, value)
}

func deleteEventFunc0(key string, value string) {
	fmt.Printf("Delete:  {key: %s , value :%s} \n", key, value)
}

func errFunc0(path string, err error) {
	fmt.Printf("error:  {path: %s , err :%s} \n", path, err.Error())
}
