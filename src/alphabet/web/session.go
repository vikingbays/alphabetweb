// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/env"
	"alphabet/sessions"
	"net/http"
)

var alphaSessionFactory sessions.AlphaSessionFactory //session存储对象 ，通过InitSession初始化
//var sessionStore *sessions.CookieStore //session存储对象 ，通过InitSession初始化
var sessionId string //sessionId信息 ，通过InitSession初始化

/*
初始化Session，当前Session存储方式支持 Cookie 和 文件系统两种
例如：初始化
  InitSession("alphabet-session-id", "", []byte("something-very-secret"))

@param sessionIdOfAll  定义sessionId信息

*/
func InitSession(sessionIdOfAll string) {

	alphaSessionFactory = sessions.CreateAlphaSessionFactory(&sessions.AlphaOptions{Path: "/", MaxAge: env.Env_Web_Sessionmaxage})

	sessionId = sessionIdOfAll
}

/*
创建一个Session对象

@param w

@param r

@return (*SessionClassInfo)

*/
func GetSession(w http.ResponseWriter, r *http.Request) sessions.AlphaSession {
	return alphaSessionFactory.GetSession(w, r)
}
