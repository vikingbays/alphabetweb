// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"container/list"
	"fmt"
	"testing"
)

//go test -v -test.run Test_List
func Test_List(t *testing.T) {

	l := list.New()
	for i := 0; i < 10; i++ {
		l.PushBack(fmt.Sprintf("data_%d", i))
	}

	for e := l.Front(); e != nil; e = e.Next() {
		t.Log("::::>>>>", e)
	}

	for e := l.Front(); e != nil; e = e.Next() {
		t.Log("---->>>>", e)
	}

}
