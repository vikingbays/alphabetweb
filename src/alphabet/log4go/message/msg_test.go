// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

import (
	"alphabet/env"
	"fmt"
	"testing"
	"time"
)

//go test -v -test.run Test_env_lang
func Test_env_lang(t *testing.T) {

	//	abtest.InitTestEnv()
	//	abtest.InitTestLog()

	env.Env_Web_I18n = "zh"
	fmt.Println(ZINF_A1_001)
	//log4go.InfoLog(INF_A1_001, 23, "ssssss")

	env.Env_Web_I18n = "en"

	//log4go.InfoLog(INF_A1_002)

	time.Sleep(time.Second)
}
