// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

////////////////////////////////////////////

//用于监控Action方法的运行时长，以及是否抛出action异常。
//
//如果需要监控，需要做两件事情
//  1、实现ActionMonitorExtends接口
//  2、实例化 ActionMonitorExtendsObject 对象，例如：
//       ActionMonitorExtendsObject = SimpleObj
type ActionMonitorExtends interface {
	AddActionMonitor(actionUrl string, appname string, startTime int64, endTime int64, err error)
}

var ActionMonitorExtendsObject ActionMonitorExtends = nil
