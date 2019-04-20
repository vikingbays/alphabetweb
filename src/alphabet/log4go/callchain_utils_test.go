// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package log4go

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

//go test -v -test.run Test_Goroutine
func Test_Goroutine(t *testing.T) {
	Goroutine(func() {
		fmt.Println("aaaaa:::", GetGID())
	})

	Goroutine(func() {
		fmt.Println("aaaaa:::", GetGID())
	})

	fmt.Println("bbbbb:::", GetGID())
}

//go test -v -test.run Test_Serialnumber
func Test_Serialnumber(t *testing.T) {
	map1 := make([]map[string]int, 1000, 1000)
	t.Log(strconv.FormatInt(time.Now().UnixNano(), 10))
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		index := i
		map1[index] = make(map[string]int)
		go func() {
			for i := 0; i < 10000; i++ {
				s1 := GetSerialnumber()
				map1[index][s1] = 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(strconv.FormatInt(time.Now().UnixNano(), 10))

	map2 := make(map[string]int)
	for _, m1 := range map1 {
		for k, v := range m1 {
			map2[k] = v
		}
	}
	t.Log(len(map2))
}

//go test -v -test.run Test_Track
func Test_Track(t *testing.T) {

	b := make([]byte, 6400)
	b = b[:runtime.Stack(b, false)]
	t.Log(string(b))
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	t.Log(n)

	t.Log(GetGID())
	t.Log(GetGID())
	t.Log(GetGID())
	t.Log(GetGID())
	t.Log(GetGID())
	t.Log(GetGID())

	var p []runtime.StackRecord = make([]runtime.StackRecord, 100, 100)
	n2, _ := runtime.GoroutineProfile(p)
	t.Log(n2)
	for i := 0; i < n2; i++ {
		p[i].Stack()
		t.Log(p[i].Stack())
	}

	t.Log(runtime.NumGoroutine())

}

//go test -v -test.run Test_AllGID
func Test_AllGID(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 2)
	}()
	go func() {
		time.Sleep(time.Second * 2)
	}()
	//GetAllID()

	t.Log(GetGID())

	t.Log(GetAllGID())

}

func GetAllID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]

	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
