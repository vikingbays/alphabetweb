// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package cache

/*
 用户触发回滚操作的错误异常
 当执行 panic(&RollbackError("sql error....."))
 就会在 EndTransAndClose中触发 rollback操作并释放连接。
*/
type RollbackError struct {
	Err error
}

func (e *RollbackError) Error() string {
	return e.Err.Error()
}

/*
 连接异常的错误
*/
type ConnectionError struct {
	Err error
}

func (e *ConnectionError) Error() string {
	return e.Err.Error()
}
