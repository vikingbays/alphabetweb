// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

/*
 定义微服务的常量信息
*/
package utils

// 如果需要重试，停顿多少毫秒.默认每次都会叠加翻倍，在Pool.Get() 方法中使用。（单位：毫秒） 1000 表示1秒
var UTILS_POOL_TRY_TIMES_PER_SLEEPTIME int = 200
