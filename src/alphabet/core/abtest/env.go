// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package abtest

import (
	"alphabet/env"
	"runtime"
	"strings"
)

/**
 * 初始化用于测试类 使用的的环境定义
 */
func InitTestEnv(projectPath string) {
	//projectPath := os.Getenv("GOLANG_PROJECT_TESTUNIT")
	if !flag_initTestEnv {
		if projectPath == "" {
			projectPath = getProjectPath()
		}
		env.InitWithoutLoad(projectPath)
		flag_initTestEnv = true
	}
}

func getProjectPath() string {
	_, filepath, _, _ := runtime.Caller(0)

	pos2 := strings.Index(filepath, "/src/alphabet/core/abtest")
	//	pos2 = pos2 + 19
	projectPath := filepath[0:pos2]
	return projectPath
}

var flag_initTestEnv bool = false
