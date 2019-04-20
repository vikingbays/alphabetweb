// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"alphabet/core/toml"
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go"
	"alphabet/log4go/message"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Client_Config_List_Toml struct {
	Client []ClientMeta_Config_List_Toml `json:"client"`
}

type ClientMeta_Config_List_Toml struct {
	GroupIds    []string                    `json:"groupIds"`
	MaxPoolSize int                         `json:"maxPoolSize"`
	ReqPerConn  int                         `json:"reqPerConn"`
	Register    []Register_Config_List_Toml `json:"register"`
}

type Server_Config_List_Toml struct {
	Server []ServerMeta_Config_List_Toml `json:"server"`
}

type ServerMeta_Config_List_Toml struct {
	GroupId        string                      `json:"groupId"`
	GroupName      string                      `json:"groupName"`
	Protocol       string                      `json:"protocol"`
	IpType         string                      `json:"ipType"`
	Ip             string                      `json:"ip"`
	Port           int                         `json:"port"`
	WebContext     string                      `json:"webContext"`
	Register       []Register_Config_List_Toml `json:"register"`
	Rpcs           []Rpc_Config_List_Toml      `json:"rpcs"`
	TicketDuration int                         `json:"ticketDuration"`
}

type Register_Config_List_Toml struct {
	Type      string   `json:"type"`
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timeout   int      `json:"timeout"`
	Root      string   `json:"root"`
}

type Rpc_Config_List_Toml struct {
	RpcId     string `json:"rpcId"`
	Url       string `json:"url"`
	Desc      string `json:"desc"`
	Available bool   `json:"available"`
}

// 初始化service配置信息
func initServAndClient() {
	serverGlobalMap = make(map[string]*Server_Info)
	clientsGlobalMap = make(map[string]*Clients_Info)

	for _, rootAppname := range env.Env_Project_Resource_Apps {
		pathOfApps := env.Env_Project_Resource + "/" + rootAppname
		if env.GetSysconfigPath() != "" {
			sysconfigPath := env.GetSysconfigPath()
			if strings.HasSuffix(sysconfigPath, "/") {
				sysconfigPath = sysconfigPath[:len(sysconfigPath)-1]
			}
			pathOfApps = pathOfApps + "/" + sysconfigPath
		}
		files, _ := ioutil.ReadDir(pathOfApps)

		serverFlag := false
		clientFlag := false
		for _, file := range files {
			if file.IsDir() {
				continue
			} else {
				if file.Name() == env.Env_Project_Resource_MS_ServerConfig_Name {
					serverFlag = true
				} else if strings.HasPrefix(file.Name(), env.Env_Project_Resource_MS_ClientConfig_Prefix) {
					clientFlag = true
				}
			}
		}

		if serverFlag {
			serverGlobalMap[rootAppname] = new(Server_Info)
			serverGlobalMap[rootAppname].rootAppname = rootAppname
		}

		if clientFlag {
			clientsGlobalMap[rootAppname] = new(Clients_Info)
			clientsGlobalMap[rootAppname].rootAppname = rootAppname
			clientsGlobalMap[rootAppname].configMap = make(map[string]Client_Config_List_Toml)
			clientsGlobalMap[rootAppname].managerMap = make(map[string]IServiceManager)
			clientsGlobalMap[rootAppname].dataMap = make(map[string]*ServiceManagerGroup)
			clientsGlobalMap[rootAppname].groupidToFilenameMap = make(map[string]string)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			} else {
				if file.Name() == env.Env_Project_Resource_MS_ServerConfig_Name {
					// 有server端的服务需要注册到注册中心
					pathconfig := pathOfApps + "/" + file.Name()
					initServerAll(pathconfig, file.Name(), rootAppname)
				} else if strings.HasPrefix(file.Name(), env.Env_Project_Resource_MS_ClientConfig_Prefix) {
					// 有client端 需要连接到注册中心
					pathconfig := pathOfApps + "/" + file.Name()
					initClientAll(pathconfig, file.Name(), rootAppname)
				}
			}
		}
	}

}

func initServerAll(pathconfig string, filename string, rootAppname string) {
	server := serverGlobalMap[rootAppname]
	server.config = initForServerConfig(pathconfig, rootAppname)

	server.rpcUrlMap = make(map[string]int)
	for _, rpc0 := range server.config.Server[0].Rpcs {
		server.rpcUrlMap[rpc0.Url] = 1
	}

	if server.config.Server[0].TicketDuration > 0 {
		SERVICE_MS_SEND_ADDR_TICKET_DURATION = time.Duration(server.config.Server[0].TicketDuration)
	}

	// 初始化，并完成注册
	if server.config.Server[0].Register[0].Type == "etcd" {
		etcdServiceManager := &EtcdServiceManager{endpoints: server.config.Server[0].Register[0].Endpoints,
			username:           server.config.Server[0].Register[0].Username,
			password:           server.config.Server[0].Register[0].Password,
			root:               server.config.Server[0].Register[0].Root,
			timeout:            time.Duration(server.config.Server[0].Register[0].Timeout) * time.Second,
			sendTicketDuration: SERVICE_MS_SEND_ADDR_TICKET_DURATION * time.Minute}
		etcdServiceManager.i = etcdServiceManager
		etcdServiceManager.typeOfService = SERVICE_MS_BASE_TYPE_SERVER
		etcdServiceManager.configName = filename
		server.manager = etcdServiceManager
	}
	if server.manager == nil {
		log4go.ErrorLog(message.ERR_MSSV_39014)
	} else {

		serviceManagerGroup := initServiceManagerGroupObject(&server.config, rootAppname)

		server.manager.RegGroupAndRpc(serviceManagerGroup)

		for _, addr := range serviceManagerGroup.Addrs {
			server.manager.RegAddr(serviceManagerGroup.GroupId, addr, rootAppname)
			break
		}

		b, _ := json.Marshal(server.config)
		//fmt.Println(string(b))
		log4go.InfoLog(message.INF_MSSV_09002, string(b))
	}

}

func initClientAll(pathconfig string, filename string, rootAppname string) {
	clients := clientsGlobalMap[rootAppname]
	clients.configMap[filename] = initForClientConfig(pathconfig)
	if clients.configMap[filename].Client[0].Register[0].Type == "etcd" {
		etcdClientManager := &EtcdServiceManager{endpoints: clients.configMap[filename].Client[0].Register[0].Endpoints,
			username: clients.configMap[filename].Client[0].Register[0].Username,
			password: clients.configMap[filename].Client[0].Register[0].Password,
			root:     clients.configMap[filename].Client[0].Register[0].Root,
			timeout:  time.Duration(clients.configMap[filename].Client[0].Register[0].Timeout) * time.Second}
		etcdClientManager.i = etcdClientManager
		etcdClientManager.typeOfService = SERVICE_MS_BASE_TYPE_CLIENT
		etcdClientManager.configName = filename
		clients.managerMap[filename] = etcdClientManager
		for _, gid := range clients.configMap[filename].Client[0].GroupIds {
			clients.groupidToFilenameMap[gid] = filename
		}
		for _, gid := range clients.configMap[filename].Client[0].GroupIds {
			clients.managerMap[filename].UpdateClientGroup(gid, false)
		}
		b, _ := json.Marshal(clients.configMap)
		log4go.InfoLog(message.INF_MSSV_09003, string(b))
		etcdClientManager.watchByConsumer()
	}
}

// 初始化ServiceManagerGroup对象
func initServiceManagerGroupObject(serverConfigListToml *Server_Config_List_Toml, rootAppname string) *ServiceManagerGroup {

	serviceManagerGroup := &ServiceManagerGroup{}

	{ // 初始化：serviceManagerGroup
		serviceManagerGroup.Addrs = make(map[string]*ServiceManagerGroupAddr)
		serviceManagerGroup.Rpcs = make(map[string]*ServiceManagerGroupRpc)

		serviceManagerGroup.GroupId = serverConfigListToml.Server[0].GroupId

		t1 := time.Now()
		serviceManagerGroup.Timestamp = fmt.Sprintf("%s%09d", t1.Format("20060102150405"), t1.UnixNano()%1000000000) //存储数据时间到nanoTime
	}

	{ // 初始化serviceManagerGroup.Addrs
		serviceManagerGroupAddr := &ServiceManagerGroupAddr{}
		//		serviceManagerGroupAddr.Group = serviceManagerGroup

		serviceManagerGroupAddr.IpType = strings.ToLower(serverConfigListToml.Server[0].IpType)

		serviceManagerGroupAddr.Timeout = SERVICE_MS_SERVER_ADDR_TIMEOUT // 10秒

		// 不区分大小写
		if strings.EqualFold(serverConfigListToml.Server[0].Protocol, "rpc_unix") {
			serviceManagerGroupAddr.Protocol = PROTOCOL_RPC_UNIX
		} else if strings.EqualFold(serverConfigListToml.Server[0].Protocol, "rpc_tcp") {
			serviceManagerGroupAddr.Protocol = PROTOCOL_RPC_TCP
		} else if strings.EqualFold(serverConfigListToml.Server[0].Protocol, "rpc_tcp_ssl") {
			serviceManagerGroupAddr.Protocol = PROTOCOL_RPC_TCP_SSL
		}

		//设置ip信息

		serviceManagerGroupAddr.Ip = serverConfigListToml.Server[0].Ip

		if serviceManagerGroupAddr.Ip == "*" {
			serviceManagerGroupAddr.Ip = getCurrentIp(serviceManagerGroupAddr.IpType)
		}

		serviceManagerGroupAddr.Port = serverConfigListToml.Server[0].Port

		serviceManagerGroupAddr.Weight = SERVICE_MS_SERVER_ADDR_WEIGHT

		//设置票据

		serviceManagerGroupAddr.Ticket = generateTicket() // 设置票据信息，根据随机数设置20位的票据信息

		serverGlobalMap[rootAppname].ticket = serviceManagerGroupAddr.Ticket
		serviceManagerGroupAddr.WebContext = serverConfigListToml.Server[0].WebContext

		serviceManagerGroup.Addrs[fmt.Sprintf("%s_%d", serviceManagerGroupAddr.Ip, serviceManagerGroupAddr.Port)] = serviceManagerGroupAddr

	}

	{ // 初始化serviceManagerGroup.Rpcs
		for _, rpc := range serverConfigListToml.Server[0].Rpcs {
			serviceManagerGroupRpc := &ServiceManagerGroupRpc{}
			//			serviceManagerGroupRpc.Group = serviceManagerGroup

			serviceManagerGroupRpc.RpcId = rpc.RpcId
			serviceManagerGroupRpc.Desc = rpc.Desc
			serviceManagerGroupRpc.Url = rpc.Url
			serviceManagerGroupRpc.Available = rpc.Available
			serviceManagerGroup.Rpcs[rpc.RpcId] = serviceManagerGroupRpc

		}

	}

	return serviceManagerGroup
}

// 初始化服务端配置文件
func initForServerConfig(pathconfig string, rootAppname string) Server_Config_List_Toml {
	var serverConfigListToml Server_Config_List_Toml
	if utils.ExistFile(pathconfig) {
		if _, err := toml.DecodeFile(pathconfig, &serverConfigListToml); err != nil {
			log4go.ErrorLog(message.ERR_MSSV_39015, err.Error())
		} else {
			if serverConfigListToml.Server[0].WebContext == "" {
				serverConfigListToml.Server[0].WebContext = env.Env_Web_Context[rootAppname]
			}
			serverConfigListToml.Server[0].Ip = strings.Replace(serverConfigListToml.Server[0].Ip,
				"${project}", env.Env_Project_Root, -1)
		}
	} else {
		log4go.ErrorLog(message.ERR_MSSV_39015, pathconfig)
	}
	return serverConfigListToml
}

// 初始化客户端配置文件
func initForClientConfig(pathconfig string) Client_Config_List_Toml {
	var clientConfigListToml Client_Config_List_Toml
	if utils.ExistFile(pathconfig) {
		if _, err := toml.DecodeFile(pathconfig, &clientConfigListToml); err != nil {
			log4go.ErrorLog(message.ERR_MSSV_39016, err.Error())
		}
	} else {
		log4go.InfoLog(message.ERR_MSSV_39016, pathconfig)
	}
	return clientConfigListToml
}
