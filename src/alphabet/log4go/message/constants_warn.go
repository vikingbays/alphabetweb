// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

/*
[WAR_CORE_69001]: 连接池获取对象失败，需要重试。
*/
var WAR_CORE_69001 MessageType = MessageType{Id: "WAR_CORE_69001",
	Cn: "在连接池基数类中，从连接池没有可以使用的对象，可能是连接池对象已经用完，也可能连接池创建有异常。在1秒后重试。",
	En: "No objects available , length of objectPool is Zero, try again after 1000 ms ."}

/*
  [WAR_CORE_69002]: 连接池获取了无效对象（判断对象是否有效，基于 ObjectFactory.Valid()验证），需要重试。
*/
var WAR_CORE_69002 MessageType = MessageType{Id: "WAR_CORE_69002",
	Cn: "在连接池基数类中，从连接池获取对象无效（判断对象是否有效，基于 ObjectFactory.Valid()验证）。在1秒后重试。",
	En: "current object is invalid when executing ObjectFactory.Valid() , try again ."}

/*
  [WAR_MSSV_69003]: 微服务客户端监听注册中心时，被退出。
*/
var WAR_MSSV_69003 MessageType = MessageType{Id: "WAR_MSSV_69003",
	Cn: "微服务客户端监听注册中心时，被退出。监听路径：%s",
	En: "When the micro service client listens to the registry, it is withdrawn. the path of monitor : %s"}

/*
  [WAR_WEB0_69004]: 微服务客户端监听注册中心时，被退出。
*/
var WAR_WEB0_69004 MessageType = MessageType{Id: "WAR_WEB0_69004",
	Cn: "在gml模板中，执行 IncludeHTML/IncludeText 的参数不存在。该参数urlStr：%s",
	En: "When executing IncludeHTML/IncludeText in gml template , this param is not exist. this param(urlStr) : %s"}

/*
  [WAR_WEB0_69005]: 微服务客户端监听注册中心时，被退出。
*/
var WAR_WEB0_69005 MessageType = MessageType{Id: "WAR_WEB0_69005",
	Cn: "在gml模板中，执行 IncludeHTML/IncludeText 的参数是外部连接，不支持。该参数urlStr：%s",
	En: "When executing IncludeHTML/IncludeText in gml template , it does not support because this param is external url . this param(urlStr) : %s"}

/*
  [WAR_MSSV_69006]: 微服务客户端监听注册中心时，被退出。
*/
var WAR_MSSV_69006 MessageType = MessageType{Id: "WAR_MSSV_69006",
	Cn: "%s 因为配置的微服务不是 rpc_tcp_ssl ,所以当前的配置的MaxPoolSize=%d, ReqPerConn=%d , 实际调整成：MaxPoolSize=%d, ReqPerConn=%d ",
	En: "%s Because protocol of MicroService is not rpc_tcp_ssl , config is MaxPoolSize=%d, ReqPerConn=%d , but now using MaxPoolSize=%d, ReqPerConn=%d ."}
