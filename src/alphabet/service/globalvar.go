// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

type Clients_Info struct {
	configMap            map[string]Client_Config_List_Toml // 配置文件信息，其中key：配置文件名称
	managerMap           map[string]IServiceManager         // 操作处理类，其中key：配置文件名称
	dataMap              map[string]*ServiceManagerGroup    //获取的微服务数据，其中key：groupId ，dataMap把多个注册中心数据整个到一起。
	groupidToFilenameMap map[string]string                  // 通过groupid，找到对应的配置文件名称，其中key：groupid，value：配置文件名称
	rootAppname          string                             // src/的根应用名
}

type Server_Info struct {
	config        Server_Config_List_Toml // 配置文件信息
	manager       IServiceManager         // 操作处理类
	ticket        string                  // 当前票据
	ticketForLast string                  // 上一个票据
	rpcUrlMap     map[string]int          // 记录rpc.Url的信息，用于验证
	rootAppname   string                  // src/的根应用名
}

//var clients *Clients_Info = new(Clients_Info)
//var server *Server_Info = new(Server_Info)

var clientsGlobalMap map[string]*Clients_Info = nil
var serverGlobalMap map[string]*Server_Info = nil

var flexServicePoolGroup *FlexServicePoolGroup = new(FlexServicePoolGroup)

// 获取server对象，可操作服务注册等操作
func GetServer(rootAppname string) *Server_Info {
	return serverGlobalMap[rootAppname]
}

// 获取server对象，可操作服务注册等操作
func GetClients(rootAppname string) *Clients_Info {
	return clientsGlobalMap[rootAppname]
}

// 根据groupId，rpcId 获取service服务信息,获取信息主要两部分，rpc信息和addr地址信息
func GetServiceInfo(groupId string, rpcId string, rootAppname string) (*ServiceManagerGroupRpc, map[string]*ServiceManagerGroupAddr) {
	serviceManagerGroupRpc := clientsGlobalMap[rootAppname].dataMap[groupId].Rpcs[rpcId]
	addrs := clientsGlobalMap[rootAppname].dataMap[groupId].Addrs
	return serviceManagerGroupRpc, addrs
}
