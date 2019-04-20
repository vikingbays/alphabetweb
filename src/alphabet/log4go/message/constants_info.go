// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

/*
MSSV --  微服务服务端
MSCI --  微服务客户端
*/

/*
  [INF_MSSV_09002]: 微服务的服务端初始化完成
*/
var INF_MSSV_09002 MessageType = MessageType{Id: "INF_MSSV_09002",
	Cn: "微服务的服务端初始化完成，创建后的 server.config 对象 是：%s",
	En: "The server of MicroService is started. server.config= %s "}

/*
  [INF_MSSV_09003]: 微服务的客户端初始化完成
*/
var INF_MSSV_09003 MessageType = MessageType{Id: "INF_MSSV_09003",
	Cn: "微服务的客户端初始化完成，创建后的 clients.configMap 对象 是：%s",
	En: "The client of MicroService is started. clients.configMap = %s "}

/*
  [INF_MSSV_09004]: 微服务客户端 从 注册中心 更新 Group配置
*/
var INF_MSSV_09004 MessageType = MessageType{Id: "INF_MSSV_09004",
	Cn: "微服务客户端 从 注册中心 更新 Group配置 ， group信息：%s",
	En: "The MicroService client updates the Group configuration from the registry . group = %s "}

/*
  [INF_MSSV_09005]: 更新微服务客户端的连接池对象
*/
var INF_MSSV_09005 MessageType = MessageType{Id: "INF_MSSV_09005",
	Cn: "更新微服务客户端的连接池对象 ， groupId=%s , addr=%s:%d",
	En: "Update the connection pool object of the MicroService client . groupId=%s , addr=%s:%d "}

/*
  [INF_MSSV_09007]: 微服务注册中心监听到的事件代码 11 ：需要修改 ticket ,只更新该地址的票据（增量模式）
*/
var INF_MSSV_09007 MessageType = MessageType{Id: "INF_MSSV_09007",
	Cn: "微服务注册中心监听到的事件代码 11 ：需要修改 ticket ,只更新该地址的票据（增量模式）， k1:%s , data01: %v",
	En: "The MicroService registry monitoring . The event code is 11 , Need to modify the ticket that updates only the address (Incremental mode) ， k1:%s , data01: %v "}

/*
  [INF_MSSV_09008]: 微服务注册中心监听到的事件代码 12 ：需要修改 timestamp ,全量更新group（全量模式）
*/
var INF_MSSV_09008 MessageType = MessageType{Id: "INF_MSSV_09008",
	Cn: "微服务注册中心监听到的事件代码 12 ：需要修改 timestamp ,全量更新group（全量模式）， k1:%s , data01: %v",
	En: "The MicroService registry monitoring . The event code is 12 , Need to modify the timestamp that updates group (Full mode) ， k1:%s , data01: %v "}

/*
  [INF_MSSV_09009]: 微服务注册中心监听到的事件代码 13 ：只新增地址 （增量模式）
*/
var INF_MSSV_09009 MessageType = MessageType{Id: "INF_MSSV_09009",
	Cn: "微服务注册中心监听到的事件代码 13 ：只新增地址 （增量模式）",
	En: "The MicroService registry monitoring . The event code is 13 , Need to add a new address (Incremental mode)  "}

/*
  [INF_MSSV_09010]: 微服务注册中心监听到的事件代码 22 ：删除了 timestamp ,本地删除group（全量模式）
*/
var INF_MSSV_09010 MessageType = MessageType{Id: "INF_MSSV_09010",
	Cn: "微服务注册中心监听到的事件代码 22 ：删除了 timestamp ,本地删除group（全量模式）,k1:%s , data01: %v",
	En: "The MicroService registry monitoring . The event code is 22 , Need to delete timestamp and delete local group (Full mode) ,k1:%s , data01: %v "}

/*
  [INF_MSSV_09011]: 微服务注册中心监听到的事件代码 23 ：删除了 timestamp ,本地删除group（全量模式）
*/
var INF_MSSV_09011 MessageType = MessageType{Id: "INF_MSSV_09011",
	Cn: "微服务注册中心监听到的事件代码 22 ：只删除本地地址，不需要请求注册中心查询（全量模式）,k1:%s , data01: %v",
	En: "The MicroService registry monitoring . The event code is 23 , Need to delete local address (Full mode) ,k1:%s , data01: %v "}

/*
  [INF_MSSV_09012]: 微服务注册中心监听处理完成后
*/
var INF_MSSV_09012 MessageType = MessageType{Id: "INF_MSSV_09012",
	Cn: "微服务注册中心监听处理完成后， 该groupId(=%s) 对应的数据 ：clients.dataMap[%s]=%s ",
	En: "Processing of the MicroService registry monitoring is complete .  groupId(=%s) , clients.dataMap[%s]=%s  "}

/*
  [INF_INIT_09013]: 应用程序，初始化，查看工程路径
*/
var INF_INIT_09013 MessageType = MessageType{Id: "INF_INIT_09013",
	Cn: "应用程序，初始化。工程路径：%s ",
	En: "Web application is initializing . project path is : %s "}

/*
  [INF_CMD0_09014]: 应用程序，初始化，查看工程路径
*/
var INF_CMD0_09014 MessageType = MessageType{Id: "INF_CMD0_09014",
	Cn: "当前的pid：%d ",
	En: "Current pid is %d "}

/*
  [INF_CMD0_09015]: 应用程序，初始化，查看工程路径
*/
var INF_CMD0_09015 MessageType = MessageType{Id: "INF_CMD0_09015",
	Cn: "生成main文件路径：%s ",
	En: "The path of main file is  %s "}

/*
  [INF_MSSV_09016]: 应用程序，初始化，查看工程路径
*/
var INF_MSSV_09016 MessageType = MessageType{Id: "INF_MSSV_09016",
	Cn: "已创建RPC 连接池(groupId=%s) ，并发请求数：%d ，创建网络连接数： %d . 连接信息：protocol=%s , addr=%s . ",
	En: "In RPC pool(groupId=%s) ,the number of Requests is %d , the number of connection object is %d . Conn Info: protocol=%s , addr=%s ."}

/*
  [INF_SQLR_09017]: 应用程序，初始化，查看工程路径
*/
var INF_SQLR_09017 MessageType = MessageType{Id: "INF_SQLR_09017",
	Cn: "已创建Database连接池(name=%s , dbType=%s) 并发请求数/创建连接对象数： %d . 连接信息： driver=%s .",
	En: "In db pool(name=%s , dbType=%s) , the number of Requests(/connection object)  is %d . Conn Info: driver=%s . "}

/*
  [INF_CACH_09018]: 应用程序，初始化，查看工程路径
*/
var INF_CACH_09018 MessageType = MessageType{Id: "INF_CACH_09018",
	Cn: "已创建Redis Cache连接池(name=%s) 并发请求数/创建连接对象数： %d . 连接信息： dataSourceName=%s .",
	En: "In redis pool(name=%s ) , the number of Requests(/connection object)  is %d . Conn Info: dataSourceName=%s . "}

/*
  [INF_MSSV_09016]: 应用程序，初始化，查看工程路径
*/
var INF_MSSV_09019 MessageType = MessageType{Id: "INF_MSSV_09019",
	Cn: "通过rpc发送请求(client -> server )，%s , Content-Type: [%s] , method: [%s] , url : [%s] , params : [%s] . ",
	En: "Request datas with rpc (client -> server )，%s , Content-Type: [%s] , method: [%s] , url : [%s] , params : [%s] . "}

//aaaaaa  INF_A1_001
var ZINF_A1_001 MessageType = MessageType{Id: "INF_A1_001",
	Cn: "中文 %d %s",
	En: "chinense"}
