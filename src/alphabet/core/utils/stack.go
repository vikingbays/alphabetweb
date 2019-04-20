// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"container/list"
)

/*
一个队列的管理，数据从队列尾插入，从队列头取出。

示例代码：
  stack:=NewStack()
  stack.Push(234)      // 插入 234
  stack.Push(666)      // 插入 666
  stack.Push(888)      // 插入 888
  stack.Pop()          // 取出 234
*/
type Stack struct {
	objectList *list.List
}

/*
初始化一个队列   stack:=NewStack()
*/
func NewStack() *Stack {
	objectList := list.New()
	return &Stack{objectList: objectList}
}

/*
插入数据，放到队列末尾

@param value  插入的数据
*/
func (stack *Stack) Push(value interface{}) {
	stack.objectList.PushBack(value)
}

/*
删除数据，从队列头删除数据

@return interface{}  删除的数据
*/
func (stack *Stack) Pop() interface{} {
	e := stack.objectList.Front()
	if e != nil {
		stack.objectList.Remove(e)
		return e.Value
	}
	return nil
}

/*
获取队列长度

@return int  队列长度
*/
func (stack *Stack) Len() int {
	return stack.objectList.Len()
}

/*
判断队列是否为空

@return bool  如果true表示队列为空。
*/
func (stack *Stack) Empty() bool {
	return stack.objectList.Len() == 0
}
