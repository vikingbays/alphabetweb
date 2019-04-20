// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package alphabet

import (
	"alphabet/cache"
	"alphabet/env"
	"alphabet/i18n"
	"alphabet/log4go"
	"alphabet/log4go/log4gohttp"
	"alphabet/log4go/message"
	"alphabet/service"
	"alphabet/sqler"
	"alphabet/web"
)

/*
完成webapp应用的初始化工作，如果设置openSqler=true，表示初始化包含数据库连接初始化。

@param projectPath      webapp工程路径。

@param openSqler        是否进行数据库连接初始化。

*/
func Init(projectPath string, openSqler bool) {
	env.Init(projectPath)
	log4go.Init()
	log4gohttp.Init()
	if env.Switch_CallChain {
		log4go.InitCallChain()
	}
	i18n.Init()

	web.Init()
	log4go.InfoLog(message.INF_INIT_09013, projectPath)
	if openSqler {
		sqler.Init()
		cache.Init()
	}
	service.Init()
}

func InitGenApplication(projectPath string) {
	env.Init(projectPath)
	log4go.Init()
	i18n.Init()
	web.InitBase()

}
