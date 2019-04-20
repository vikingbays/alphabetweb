// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package abtest

import (
	"alphabet/env"
	"alphabet/log4go"
	"testing"
)

/**
 * 初始化用于测试类 使用的的日志类。
 */
func InitTestLog(logconfigFile string) {
	// 默认日志调用链关闭
	if !flag_initTestLog {
		env.Switch_CallChain = false
		if logconfigFile == "" {
			logconfigFile = env.Env_Project_Resource + "/alphabet/core/abtest/logconfig_test.toml"
		}
		log4go.LINENO_CALLER_LEVEL = 3
		log4go.InitLoggerInstance(log4go.InitLogger(logconfigFile, env.Env_Project_Root))
		flag_initTestLog = true
	}
}

// 标记是否被初始化
var flag_initTestLog bool = false

/*
* 用于测试类的日志输出，如果加载了log4go，那么就使用log4go输出。如果没有加载了log4go，就使用 t
 */
func InfoTestLog(t *testing.T, arg0 string, args ...interface{}) {
	if log4go.IsInitCompleted() { // 判断log4go是否初始化，如果初始化，就使用log4go输出
		log4go.GetL().Info(3, arg0, args...)
	} else { // 如果log4go没有初始化，就使用 t.Log()
		t.Logf(arg0, args...)
	}
}

/*
* 用于测试类的日志输出，如果加载了log4go，那么就使用log4go输出。如果没有加载了log4go，就使用 t
* 如果使用了log4go，那么会给错误日志打标记
 */
func ErrorTestLog(t *testing.T, arg0 string, args ...interface{}) {
	if log4go.IsInitCompleted() { // 判断log4go是否初始化，如果初始化，就使用log4go输出
		log4go.GetL().Error(3, arg0, args...)
		t.Fail()
	} else { // 如果log4go没有初始化，就使用 t.Error()
		t.Errorf(arg0, args...)
	}
}

/*
* 用于测试类的日志输出，如果加载了log4go，那么就使用log4go输出。如果没有加载了log4go，就使用 t
 */
func DebugTestLog(t *testing.T, arg0 string, args ...interface{}) {
	if log4go.IsInitCompleted() { // 判断log4go是否初始化，如果初始化，就使用log4go输出
		log4go.GetL().Debug(3, arg0, args...)
	} else { // 如果log4go没有初始化，就使用 t.Log()
		t.Logf(arg0, args...)
	}
}
