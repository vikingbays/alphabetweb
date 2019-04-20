// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package message

/*
[DEG_CORE_79001]: 当前etcd的连接状态
*/
var DEG_CORE_79001 MessageType = MessageType{Id: "DEG_CORE_79001",
	Cn: "当前etcd的连接状态；%s",
	En: "The state of etcdConnection is: %s"}

/*
  [DEG_CORE_79002]: 无效的数据】当使用etcd的批量更新时，有部分key被重复更新
*/
var DEG_CORE_79002 MessageType = MessageType{Id: "DEG_CORE_79002",
	Cn: "【无效的数据】当使用etcd的批量更新时，有部分key被重复更新，检测到被重复更新的key数据： { key : %s , value: %s , oper: %d,resp: %v }",
	En: "（invalid data）when batch update in etcd , this key is duplicate : { key : %s , value: %s , oper: %d,resp: %v } "}

/*
   [DEG_CORE_79003]: 【有效的数据】当使用etcd的批量更新时，已经被更新的数据
*/
var DEG_CORE_79003 MessageType = MessageType{Id: "DEG_CORE_79003",
	Cn: "【有效的数据】当使用etcd的批量更新时，已经被更新的数据： { key : %s , value: %s , oper: %d,resp: %v }",
	En: "（valid data）when batch update in etcd , this key is updated  : { key : %s , value: %s , oper: %d,resp: %v } "}

/*
   [DEG_CORE_79004]: 在微服务连接池中，获取到的rpc对象
*/
var DEG_CORE_79004 MessageType = MessageType{Id: "DEG_CORE_79004",
	Cn: "在微服务连接池中，获取到的rpc对象是：poolGroupLength=%d , currNum=%d , groupId=%s , IpAndPort=%s ",
	En: "In the MicroService connection pool, the RPC object is obtained.  : poolGroupLength=%d , currNum=%d , groupId=%s , IpAndPort=%s "}

/*
   [DEG_MSSV_79005]: 在微服务连接池中，获取到的rpc对象
*/
var DEG_MSSV_79005 MessageType = MessageType{Id: "DEG_MSSV_79005",
	Cn: "微服务客户端准备创建客户端连接。Protocol=%s  , ip=%s , port=%d  ",
	En: "The MicroService client will create a connection.  Protocol=%s  , ip=%s , port=%d "}

/*
   [DEG_SQLR_79006]: 本次数据库操作信息
*/
var DEG_SQLR_79006 MessageType = MessageType{Id: "DEG_SQLR_79006",
	Cn: " 本次操作信息。 appName: %s \n sqlname: %s \n template:%s \n sql:    %s \n params:  %v \n  ",
	En: "This operation is >>  appName: %s \n sqlname: %s \n template:%s \n sql:    %s \n params:  %v \n  "}

/*
   [DEG_SQLR_79007]: 数据库SQL基准测试 开始时间
*/
var DEG_SQLR_79007 MessageType = MessageType{Id: "DEG_SQLR_79007",
	Cn: "[SQL 基准测试] %s 开始 : %s ",
	En: "[SQL Performance] %s begin : %s"}

/*
   [DEG_SQLR_79008]: 数据库SQL基准测试 结束时间。
*/
var DEG_SQLR_79008 MessageType = MessageType{Id: "DEG_SQLR_79008",
	Cn: "[SQL 基准测试] %s 结束 :  %s ",
	En: "[SQL Performance] %s end :  %s"}

/*
   [DEG_SQLR_79009]: 本次数据库操作信息
*/
var DEG_SQLR_79009 MessageType = MessageType{Id: "DEG_SQLR_79009",
	Cn: " 本次操作信息。 appName: %s \n sqlname: %s  ",
	En: "This operation is >> appName: %s \n sqlname: %s  "}

/*
   [DEG_CALL_79010]: 本次数据库操作信息
*/
var DEG_CALL_79010 MessageType = MessageType{Id: "DEG_CALL_79010",
	Cn: " 当前记录的协程信息caller： %v  ",
	En: "The map of current gids(caller) : %v "}

/*
	[DEG_MSSV_79011]: 记录微服务的连接池对象的创建信息，接收参数：groupId 和 addr
*/
var DEG_MSSV_79011 MessageType = MessageType{Id: "DEG_MSSV_79011",
	Cn: "在微服务连接池中创建一个Rpc对象，其中 GroupId=%s, Addr=%s",
	En: "Create a rpcObject in servicePools ， GroupId=%s, Addr=%s"}

/*
  [DEG_SQLR_79012]: 正在创建一个数据库连接 。
*/
var DEG_SQLR_79012 MessageType = MessageType{Id: "DEG_SQLR_79012",
	Cn: "正在创建一个数据库连接 (正在进行中)。 c.DriverName=%s, c.DataSourceName=%s ",
	En: "Creating a connection of database (doing). c.DriverName=%s, c.DataSourceName=%s "}
